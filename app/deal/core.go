package deal

import (
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/seat"
)

type ID = id.ADT

type DealSpec struct {
	Name string
}

type DealRef struct {
	ID   id.ADT
	Name string
}

// aka Configuration or Eta
type DealRoot struct {
	ID       id.ADT
	Name     string
	Children []DealRef
	Seats    []seat.SeatRef
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Ident
var (
	ConvertRootToRef func(DealRoot) DealRef
)

type DealApi interface {
	Create(DealSpec) (DealRoot, error)
	Retrieve(id.ADT) (DealRoot, error)
	RetreiveAll() ([]DealRef, error)
	Establish(KinshipSpec) error
	Involve(PartSpec) (PartRoot, error)
	Take(TranSpec) error
}

type dealService struct {
	deals    dealRepo
	seats    seat.SeatApi
	parts    partRepo
	chnls    chnl.Repo
	procs    step.Repo[step.ProcRoot]
	msgs     step.Repo[step.MsgRoot]
	srvs     step.Repo[step.SrvRoot]
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newDealService(
	deals dealRepo,
	seats seat.SeatApi,
	parts partRepo,
	chnls chnl.Repo,
	procs step.Repo[step.ProcRoot],
	msgs step.Repo[step.MsgRoot],
	srvs step.Repo[step.SrvRoot],
	states state.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *dealService {
	name := slog.String("name", "dealService")
	return &dealService{deals, seats, parts, chnls, procs, msgs, srvs, states, kinships, l.With(name)}
}

func (s *dealService) Create(spec DealSpec) (DealRoot, error) {
	s.log.Debug("deal creation started", slog.Any("spec", spec))
	root := DealRoot{
		ID:   id.New(),
		Name: spec.Name,
	}
	err := s.deals.Insert(root)
	if err != nil {
		s.log.Error("deal insertion failed",
			slog.Any("reason", err),
			slog.Any("root", root),
		)
		return root, err
	}
	s.log.Debug("deal creation succeeded", slog.Any("root", root))
	return root, nil
}

func (s *dealService) Retrieve(id id.ADT) (DealRoot, error) {
	root, err := s.deals.SelectByID(id)
	if err != nil {
		return DealRoot{}, err
	}
	root.Children, err = s.deals.SelectChildren(id)
	if err != nil {
		return DealRoot{}, err
	}
	return root, nil
}

func (s *dealService) RetreiveAll() ([]DealRef, error) {
	return s.deals.SelectAll()
}

func (s *dealService) Establish(spec KinshipSpec) error {
	s.log.Debug("kinship establishment started", slog.Any("spec", spec))
	var children []DealRef
	for _, id := range spec.ChildIDs {
		children = append(children, DealRef{ID: id})
	}
	root := KinshipRoot{
		Parent:   DealRef{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinships.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("kinship establishment succeeded", slog.Any("root", root))
	return nil
}

func (s *dealService) Involve(spec PartSpec) (PartRoot, error) {
	s.log.Debug("seat involvement started", slog.Any("spec", spec))
	seat, err := s.seats.Retrieve(spec.SeatID)
	if err != nil {
		s.log.Error("seat selection failed",
			slog.Any("reason", err),
			slog.Any("spec", spec),
		)
		return PartRoot{}, err
	}
	via := chnl.Root{
		ID:   id.New(),
		Name: seat.Via.Name,
		St:   seat.Via.St,
	}
	err = s.chnls.Insert(via)
	if err != nil {
		s.log.Error("via insertion failed",
			slog.Any("reason", err),
			slog.Any("via", via),
		)
		return PartRoot{}, err
	}
	ctx := make(map[chnl.Name]chnl.ID, len(spec.Ctx))
	if len(spec.Ctx) > 0 {
		var chnls []chnl.Root
		for _, preID := range spec.Ctx {
			ch := chnl.Root{
				ID:    id.New(),
				PreID: preID,
			}
			chnls = append(chnls, ch)
		}
		chnls, err = s.chnls.InsertCtx(chnls)
		if err != nil {
			s.log.Error("ctx insertion failed",
				slog.Any("reason", err),
				slog.Any("ctx", chnls),
			)
			return PartRoot{}, err
		}
		for _, ch := range chnls {
			ctx[ch.Name] = ch.ID
		}
	}
	root := PartRoot{
		PartID: id.New(),
		DealID: spec.DealID,
		SeatID: spec.SeatID,
		PAK:    ak.New(),
		CAK:    ak.New(),
		PID:    via.ID,
		Ctx:    ctx,
	}
	err = s.parts.Insert(root)
	if err != nil {
		s.log.Error("participation insertion failed",
			slog.Any("reason", err),
			slog.Any("root", root),
		)
		return PartRoot{}, err
	}
	s.log.Debug("seat involvement succeeded", slog.Any("root", root))
	return root, nil
}

func (s *dealService) Take(spec TranSpec) error {
	if spec.Term == nil {
		panic(step.ErrUnexpectedTerm(spec.Term))
	}
	s.log.Debug("transition taking started", slog.Any("spec", spec))
	part, err := s.parts.SelectByID(spec.PartID)
	if err != nil {
		s.log.Error("participation selection failed",
			slog.Any("reason", err),
			slog.Any("id", spec.PartID),
		)
		return err
	}
	switch term := spec.Term.(type) {
	case step.CloseSpec:
		curChnl, err := s.chnls.SelectByID(term.A)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", term.A),
			)
			return err
		}
		if curChnl.St == nil {
			err = errAlreadyClosedChannel(chnl.ConvertRootToRef(curChnl))
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			s.log.Error("state selection failed",
				slog.Any("reason", err),
				slog.Any("st", curChnl.St),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(term, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(term, curSt)
		} else {
			err = fmt.Errorf("unexpected access key: %s", spec.AgentAK)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if err != nil {
			s.log.Error("type checking failed",
				slog.Any("reason", err),
				slog.Any("term", term),
			)
			return err
		}
		// TODO выборка с проверкой потребления
		srv, err := s.srvs.SelectByCh(curChnl.ID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		if srv == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				VID: term.A,
				Val: term,
			}
			err = s.msgs.Insert(newMsg)
			if err != nil {
				s.log.Error("message insertion failed",
					slog.Any("reason", err),
					slog.Any("msg", newMsg),
				)
				return err
			}
			s.log.Debug("transition taking half done", slog.Any("msg", newMsg))
			return nil
		}
		wait, ok := srv.Cont.(step.WaitSpec)
		if !ok {
			err = fmt.Errorf("unexpected cont type: want %T, got %T", step.WaitSpec{}, srv.Cont)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("cont", srv.Cont),
			)
			return err
		}
		// consume channel
		finChnl := chnl.Root{
			ID:    id.New(),
			Name:  curChnl.Name,
			PreID: curChnl.ID,
			St:    nil,
		}
		err = s.chnls.Insert(finChnl)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("chnl", finChnl),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		return s.Take(TranSpec{PartID: spec.PartID, AgentAK: spec.AgentAK, Term: wait.Cont})
	case step.WaitSpec:
		curChnl, err := s.chnls.SelectByID(term.X)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", term.X),
			)
			return err
		}
		if curChnl.St == nil {
			err = errAlreadyClosedChannel(chnl.ConvertRootToRef(curChnl))
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			s.log.Error("state selection failed",
				slog.Any("reason", err),
				slog.Any("st", curChnl.St),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(term, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(term, curSt)
		} else {
			err = fmt.Errorf("unexpected access key: %s", spec.AgentAK)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if err != nil {
			s.log.Error("type checking failed",
				slog.Any("reason", err),
				slog.Any("term", term),
			)
			return err
		}
		msg, err := s.msgs.SelectByCh(curChnl.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		if msg == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				VID:  term.X,
				Cont: term,
			}
			err = s.srvs.Insert(newSrv)
			if err != nil {
				s.log.Error("service insertion failed",
					slog.Any("reason", err),
					slog.Any("srv", newSrv),
				)
				return err
			}
			s.log.Debug("transition taking half done", slog.Any("srv", newSrv))
			return nil
		}
		_, ok := msg.Val.(step.CloseSpec)
		if !ok {
			err = fmt.Errorf("unexpected val type: want %T, got %T", step.CloseSpec{}, msg.Val)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("val", msg.Val),
			)
			return err
		}
		// consume channel
		finChnl := chnl.Root{
			ID:    id.New(),
			Name:  curChnl.Name,
			PreID: curChnl.ID,
			St:    nil,
		}
		err = s.chnls.Insert(finChnl)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("chnl", finChnl),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		return s.Take(TranSpec{PartID: spec.PartID, AgentAK: spec.AgentAK, Term: term.Cont})
	case step.SendSpec:
		curChnl, err := s.chnls.SelectByID(term.A)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", term.A),
			)
			return err
		}
		if curChnl.St == nil {
			err = errAlreadyClosedChannel(chnl.ConvertRootToRef(curChnl))
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			s.log.Error("state selection failed",
				slog.Any("reason", err),
				slog.Any("st", curChnl.St),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(term, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(term, curSt)
		} else {
			err = fmt.Errorf("unexpected access key: %s", spec.AgentAK)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if err != nil {
			s.log.Error("type checking failed",
				slog.Any("reason", err),
				slog.Any("term", term),
			)
			return err
		}
		srv, err := s.srvs.SelectByCh(curChnl.ID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		if srv == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				VID: term.A,
				Val: term,
			}
			err = s.msgs.Insert(newMsg)
			if err != nil {
				s.log.Error("message insertion failed",
					slog.Any("reason", err),
					slog.Any("msg", newMsg),
				)
				return err
			}
			s.log.Debug("transition taking half done", slog.Any("msg", newMsg))
			return nil
		}
		recv, ok := srv.Cont.(step.RecvSpec)
		if !ok {
			err = fmt.Errorf("unexpected cont type: want %T, got %T", step.RecvSpec{}, srv.Cont)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("cont", srv.Cont),
			)
			return err
		}
		newChnl := chnl.Root{
			ID:    id.New(),
			Name:  curChnl.Name,
			PreID: curChnl.ID,
			St:    curSt.(state.Prod).Next(),
		}
		err = s.chnls.Insert(newChnl)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("chnl", newChnl),
			)
			return err
		}
		recv.Cont = step.Subst(recv.Cont, recv.X, newChnl.ID)
		recv.Cont = step.Subst(recv.Cont, recv.Y, term.B)
		s.log.Debug("transition taking succeeded")
		return s.Take(TranSpec{PartID: spec.PartID, AgentAK: spec.AgentAK, Term: recv.Cont})
	case step.RecvSpec:
		curChnl, err := s.chnls.SelectByID(term.X)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", term.X),
			)
			return err
		}
		if curChnl.St == nil {
			err = errAlreadyClosedChannel(chnl.ConvertRootToRef(curChnl))
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			s.log.Error("state selection failed",
				slog.Any("reason", err),
				slog.Any("st", curChnl.St),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(term, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(term, curSt)
		} else {
			err = fmt.Errorf("unexpected access key: %s", spec.AgentAK)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if err != nil {
			s.log.Error("type checking failed",
				slog.Any("reason", err),
				slog.Any("term", term),
			)
			return err
		}
		msg, err := s.msgs.SelectByCh(curChnl.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		if msg == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				VID:  term.X,
				Cont: term,
			}
			err = s.srvs.Insert(newSrv)
			if err != nil {
				s.log.Error("service insertion failed",
					slog.Any("reason", err),
					slog.Any("srv", newSrv),
				)
				return err
			}
			s.log.Debug("transition taking half done", slog.Any("srv", newSrv))
			return nil
		}
		val, ok := msg.Val.(step.SendSpec)
		if !ok {
			err = fmt.Errorf("unexpected val type: want %T, got %T", step.SendSpec{}, msg.Val)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("val", msg.Val),
			)
			return err
		}
		newChnl := chnl.Root{
			ID:    id.New(),
			Name:  curChnl.Name,
			PreID: curChnl.ID,
			St:    curSt.(state.Prod).Next(),
		}
		err = s.chnls.Insert(newChnl)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("chnl", newChnl),
			)
			return err
		}
		term.Cont = step.Subst(term.Cont, term.X, newChnl.ID)
		term.Cont = step.Subst(term.Cont, term.Y, val.B)
		s.log.Debug("transition taking succeeded")
		return s.Take(TranSpec{PartID: spec.PartID, AgentAK: spec.AgentAK, Term: term.Cont})
	case step.LabSpec:
		curChnl, err := s.chnls.SelectByID(term.C)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", term.C),
			)
			return err
		}
		if curChnl.St == nil {
			err = errAlreadyClosedChannel(chnl.ConvertRootToRef(curChnl))
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			s.log.Error("state selection failed",
				slog.Any("reason", err),
				slog.Any("st", curChnl.St),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(term, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(term, curSt)
		} else {
			err = fmt.Errorf("unexpected access key: %s", spec.AgentAK)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if err != nil {
			s.log.Error("type checking failed",
				slog.Any("reason", err),
				slog.Any("term", term),
			)
			return err
		}
		srv, err := s.srvs.SelectByCh(curChnl.ID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		if srv == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				VID: term.C,
				Val: term,
			}
			err = s.msgs.Insert(newMsg)
			if err != nil {
				s.log.Error("message insertion failed",
					slog.Any("reason", err),
					slog.Any("msg", newMsg),
				)
				return err
			}
			s.log.Debug("transition taking half done", slog.Any("msg", newMsg))
			return nil
		}
		cont, ok := srv.Cont.(step.CaseSpec)
		if !ok {
			err = fmt.Errorf("unexpected cont type: want %T, got %T", step.CaseSpec{}, srv.Cont)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("cont", srv.Cont),
			)
			return err
		}
		newChnl := chnl.Root{
			ID:    id.New(),
			Name:  curChnl.Name,
			PreID: curChnl.ID,
			St:    curSt.(state.Sum).Next(term.L),
		}
		err = s.chnls.Insert(newChnl)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("chnl", newChnl),
			)
			return err
		}
		branch := step.Subst(cont.Conts[term.L], cont.Z, newChnl.ID)
		return s.Take(TranSpec{PartID: spec.PartID, AgentAK: spec.AgentAK, Term: branch})
	case step.CaseSpec:
		curChnl, err := s.chnls.SelectByID(term.Z)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", term.Z),
			)
			return err
		}
		if curChnl.St == nil {
			err = errAlreadyClosedChannel(chnl.ConvertRootToRef(curChnl))
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			s.log.Error("state selection failed",
				slog.Any("reason", err),
				slog.Any("st", curChnl.St),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(term, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(term, curSt)
		} else {
			err = fmt.Errorf("unexpected access key: %s", spec.AgentAK)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if err != nil {
			s.log.Error("type checking failed",
				slog.Any("reason", err),
				slog.Any("term", term),
			)
			return err
		}
		msg, err := s.msgs.SelectByCh(curChnl.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		if msg == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				VID:  term.Z,
				Cont: term,
			}
			err = s.srvs.Insert(newSrv)
			if err != nil {
				s.log.Error("service insertion failed",
					slog.Any("reason", err),
					slog.Any("srv", newSrv),
				)
				return err
			}
			s.log.Debug("transition taking half done", slog.Any("srv", newSrv))
			return nil
		}
		val, ok := msg.Val.(step.LabSpec)
		if !ok {
			err = fmt.Errorf("unexpected val type: want %T, got %T", step.LabSpec{}, msg.Val)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("val", msg.Val),
			)
			return err
		}
		newChnl := chnl.Root{
			ID:    id.New(),
			Name:  curChnl.Name,
			PreID: curChnl.ID,
			St:    curSt.(state.Sum).Next(val.L),
		}
		err = s.chnls.Insert(newChnl)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("chnl", newChnl),
			)
			return err
		}
		branch := step.Subst(term.Conts[val.L], term.Z, newChnl.ID)
		return s.Take(TranSpec{spec.DealID, spec.PartID, spec.AgentAK, branch})
	case step.SpawnSpec:
		part, err := s.Involve(PartSpec{spec.DealID, term.DecID, term.Ctx})
		if err != nil {
			return err
		}
		term.Cont = step.Subst(term.Cont, term.C, part.PID)
		return s.Take(TranSpec{spec.DealID, spec.PartID, spec.AgentAK, term.Cont})
	default:
		panic(step.ErrUnexpectedTerm(spec.Term))
	}
}

type dealRepo interface {
	Insert(DealRoot) error
	SelectAll() ([]DealRef, error)
	SelectByID(ID) (DealRoot, error)
	SelectChildren(ID) ([]DealRef, error)
	SelectSeats(ID) ([]seat.SeatRef, error)
}

// Kinship Relation
type KinshipSpec struct {
	ParentID ID
	ChildIDs []ID
}

type KinshipRoot struct {
	Parent   DealRef
	Children []DealRef
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// Participation aka Spawn (but external)
type PartSpec struct {
	DealID ID
	SeatID seat.ID
	Ctx    map[chnl.Name]chnl.ID
}

type PartRoot struct {
	PartID ID
	DealID ID
	SeatID seat.ID
	// Producible Access Key
	PAK ak.ADT
	// Consumables Access Key
	CAK ak.ADT
	// Producible Channel ID
	PID chnl.ID
	// Consumable Channel IDs
	Ctx map[chnl.Name]chnl.ID
}

type partRepo interface {
	Insert(PartRoot) error
	SelectByID(ID) (PartRoot, error)
}

// Transition
type TranSpec struct {
	DealID  ID
	PartID  ID
	AgentAK ak.ADT
	Term    step.Term
}

// aka checkExp
func (s *dealService) checkProducer(t step.Term, st state.Root) error {
	switch term := t.(type) {
	case step.CloseSpec:
		return checkProducer(st, state.OneRoot{})
	case step.WaitSpec:
		return checkProducer(st, state.OneRoot{})
	case step.SendSpec:
		// check value
		want, ok := st.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.TensorRoot{}, st)
		}
		got, err := s.chnls.SelectByID(term.B)
		if err != nil {
			return err
		}
		err = checkProducer(got.St, want.B)
		if err != nil {
			return err
		}
		// no cont to check
		return nil
	case step.RecvSpec:
		want, ok := st.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.LolliRoot{}, st)
		}
		// check value
		got, err := s.chnls.SelectByID(term.Y)
		if err != nil {
			return err
		}
		err = checkProducer(got.St, want.Y)
		if err != nil {
			return err
		}
		// check cont
		return s.checkProducer(term.Cont, want.Z)
	case step.LabSpec:
		want, ok := st.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.PlusRoot{}, st)
		}
		_, ok = want.Choices[term.L]
		if !ok {
			return fmt.Errorf("label mismatch: want %q, got nothing", term.L)
		}
		// no cont to check
		return nil
	case step.CaseSpec:
		want, ok := st.(state.WithRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.WithRoot{}, st)
		}
		if len(term.Conts) != len(want.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v branches", len(want.Choices), len(term.Conts))
		}
		for wantL, wantCh := range want.Choices {
			gotBr, ok := term.Conts[wantL]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantL)
			}
			err := s.checkProducer(gotBr, wantCh)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(step.ErrUnexpectedTerm(t))
	}
}

func (s *dealService) checkConsumer(t step.Term, st state.Root) error {
	switch term := t.(type) {
	case step.CloseSpec:
		return checkConsumer(st, state.OneRoot{})
	case step.WaitSpec:
		return checkConsumer(st, state.OneRoot{})
	case step.SendSpec:
		// check value
		want, ok := st.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.LolliRoot{}, st)
		}
		got, err := s.chnls.SelectByID(term.B)
		if err != nil {
			return err
		}
		err = checkConsumer(got.St, want.Y)
		if err != nil {
			return err
		}
		// no cont to check
		return nil
	case step.RecvSpec:
		want, ok := st.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.TensorRoot{}, st)
		}
		// check value
		got, err := s.chnls.SelectByID(term.Y)
		if err != nil {
			return err
		}
		err = checkConsumer(got.St, want.B)
		if err != nil {
			return err
		}
		// check cont
		return s.checkConsumer(term.Cont, want.C)
	case step.LabSpec:
		want, ok := st.(state.WithRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.WithRoot{}, st)
		}
		_, ok = want.Choices[term.L]
		if !ok {
			return fmt.Errorf("label mismatch: want %q, got nothing", term.L)
		}
		// no cont to check
		return nil
	case step.CaseSpec:
		want, ok := st.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.PlusRoot{}, st)
		}
		if len(term.Conts) != len(want.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v branches", len(want.Choices), len(term.Conts))
		}
		for wantL, wantCh := range want.Choices {
			gotBr, ok := term.Conts[wantL]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantL)
			}
			err := s.checkConsumer(gotBr, wantCh)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(step.ErrUnexpectedTerm(t))
	}
}

// aka eqtp
func checkProducer(got, want state.Root) error {
	switch wantSt := want.(type) {
	case state.OneRef:
		_, ok := got.(state.OneRef)
		if !ok {
			return fmt.Errorf("state ref mismatch: want %T, got %#v", want, got)
		}
		return nil
	case state.OneRoot:
		_, ok := got.(state.OneRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		return nil
	case state.TensorRef:
		gotSt, ok := got.(state.TensorRef)
		if !ok {
			return fmt.Errorf("state ref mismatch: want %T, got %#v", want, got)
		}
		if gotSt != wantSt {
			return fmt.Errorf("state ref mismatch: want %#v, got %#v", want, got)
		}
		return nil
	case state.TensorRoot:
		gotSt, ok := got.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		err := checkProducer(gotSt.B, wantSt.B)
		if err != nil {
			return err
		}
		return checkProducer(gotSt.C, wantSt.C)
	case state.LolliRoot:
		gotSt, ok := got.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		err := checkProducer(gotSt.Y, wantSt.Y)
		if err != nil {
			return err
		}
		return checkProducer(gotSt.Z, wantSt.Z)
	case state.PlusRoot:
		gotSt, ok := got.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v choices", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantL, wantCh := range wantSt.Choices {
			gotCh, ok := gotSt.Choices[wantL]
			if !ok {
				return fmt.Errorf("state mismatch: want label %q, got nothing", wantL)
			}
			err := checkProducer(gotCh, wantCh)
			if err != nil {
				return err
			}
		}
		return nil
	case state.WithRoot:
		gotSt, ok := got.(state.WithRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v choices", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantL, wantCh := range wantSt.Choices {
			gotCh, ok := gotSt.Choices[wantL]
			if !ok {
				return fmt.Errorf("state mismatch: want label %q, got nothing", wantL)
			}
			err := checkProducer(gotCh, wantCh)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(state.ErrUnexpectedRoot(want))
	}
}

func checkConsumer(got, want state.Root) error {
	switch wantSt := want.(type) {
	case state.OneRef:
		_, ok := got.(state.OneRef)
		if !ok {
			return fmt.Errorf("state ref mismatch: want %T, got %#v", want, got)
		}
		return nil
	case state.OneRoot:
		_, ok := got.(state.OneRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		return nil
	case state.TensorRef:
		gotSt, ok := got.(state.TensorRef)
		if !ok {
			return fmt.Errorf("state ref mismatch: want %T, got %#v", want, got)
		}
		if gotSt != wantSt {
			return fmt.Errorf("state ref mismatch: want %#v, got %#v", want, got)
		}
		return nil
	case state.TensorRoot:
		gotSt, ok := got.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		err := checkConsumer(gotSt.Y, wantSt.B)
		if err != nil {
			return err
		}
		return checkConsumer(gotSt.Z, wantSt.C)
	case state.LolliRoot:
		gotSt, ok := got.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		err := checkConsumer(gotSt.B, wantSt.Y)
		if err != nil {
			return err
		}
		return checkConsumer(gotSt.C, wantSt.Z)
	case state.PlusRoot:
		gotSt, ok := got.(state.WithRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v choices", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantL, wantCh := range wantSt.Choices {
			gotCh, ok := gotSt.Choices[wantL]
			if !ok {
				return fmt.Errorf("state mismatch: want label %q, got nothing", wantL)
			}
			err := checkConsumer(gotCh, wantCh)
			if err != nil {
				return err
			}
		}
		return nil
	case state.WithRoot:
		gotSt, ok := got.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v choices", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantL, wantCh := range wantSt.Choices {
			gotCh, ok := gotSt.Choices[wantL]
			if !ok {
				return fmt.Errorf("state mismatch: want label %q, got nothing", wantL)
			}
			err := checkConsumer(gotCh, wantCh)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(state.ErrUnexpectedRoot(want))
	}
}

func errAlreadyClosedChannel(ref chnl.Ref) error {
	return fmt.Errorf("channel already finalized %+v", ref)
}
