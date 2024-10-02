package deal

import (
	"fmt"
	"log/slog"

	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/seat"
)

type Name = string
type ID = id.ADT

type DealSpec struct {
	Name Name
}

type DealRef struct {
	ID   ID
	Name Name
}

// aka Configuration or Eta
type DealRoot struct {
	ID       ID
	Name     Name
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
	Retrieve(ID) (DealRoot, error)
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
	dec, err := s.seats.Retrieve(spec.SeatID)
	if err != nil {
		s.log.Error("seat selection failed",
			slog.Any("reason", err),
			slog.Any("spec", spec),
		)
		return PartRoot{}, err
	}
	if len(dec.Ctx) != len(spec.Ctx) {
		err = fmt.Errorf("ctx mismatch: want %v items, got %v items", len(dec.Ctx), len(spec.Ctx))
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
			slog.Any("spec", spec),
		)
		return PartRoot{}, err
	}
	newCtx := make(map[chnl.Name]chnl.ID, len(spec.Ctx))
	if len(spec.Ctx) > 0 {
		curChnls, err := s.chnls.SelectMany(maps.Values(spec.Ctx))
		if err != nil {
			s.log.Error("ctx selection failed",
				slog.Any("reason", err),
				slog.Any("spec", spec),
			)
			return PartRoot{}, err
		}
		for i, got := range curChnls {
			// TODO обеспечить порядок
			err = checkState(got.St, dec.Ctx[i].St)
			if err != nil {
				s.log.Error("type checking failed",
					slog.Any("reason", err),
					slog.Any("spec", spec),
				)
				return PartRoot{}, err
			}
		}
		var newChnls []chnl.Root
		for _, preID := range spec.Ctx {
			ch := chnl.Root{
				ID:    id.New(),
				PreID: preID,
			}
			newChnls = append(newChnls, ch)
		}
		newChnls, err = s.chnls.InsertCtx(newChnls)
		if err != nil {
			s.log.Error("ctx insertion failed",
				slog.Any("reason", err),
				slog.Any("ctx", newChnls),
			)
			return PartRoot{}, err
		}
		for _, ch := range newChnls {
			newCtx[ch.Name] = ch.ID
		}
	}
	via := chnl.Root{
		ID:   id.New(),
		Name: dec.Via.Name,
		StID: dec.Via.StID,
		St:   dec.Via.St,
	}
	err = s.chnls.Insert(via)
	if err != nil {
		s.log.Error("via insertion failed",
			slog.Any("reason", err),
			slog.Any("via", via),
		)
		return PartRoot{}, err
	}
	root := PartRoot{
		PartID: id.New(),
		DealID: spec.DealID,
		SeatID: spec.SeatID,
		PAK:    ak.New(),
		CAK:    ak.New(),
		PID:    via.ID,
		Ctx:    newCtx,
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
	chIDs := step.CollectChnlIDs(spec.Term, []chnl.ID{})
	cfg, err := s.chnls.SelectCfg(chIDs)
	if err != nil {
		s.log.Error("cfg selection failed",
			slog.Any("reason", err),
			slog.Any("ids", chIDs),
		)
		return err
	}
	stIDs := chnl.CollectStIDs(maps.Values(cfg))
	env, err := s.states.SelectEnv(stIDs)
	if err != nil {
		s.log.Error("env selection failed",
			slog.Any("reason", err),
			slog.Any("ids", stIDs),
		)
		return err
	}
	part, err := s.parts.SelectByID(spec.PartID)
	if err != nil {
		s.log.Error("participation selection failed",
			slog.Any("reason", err),
			slog.Any("id", spec.PartID),
		)
		return err
	}
	return s.takeRecur(env, cfg, part, spec)
}

func (s *dealService) takeRecur(
	env map[state.ID]state.Root,
	cfg map[chnl.ID]chnl.Root,
	part PartRoot,
	spec TranSpec,
) (err error) {
	switch term := spec.Term.(type) {
	case step.CloseSpec:
		curChnl, ok := cfg[term.A]
		if !ok {
			err := chnl.ErrDoesNotExist(term.A)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if curChnl.St == nil {
			err := chnl.ErrAlreadyClosed
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, ok := env[curChnl.StID]
		if !ok {
			err := state.ErrDoesNotExist(curChnl.StID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(cfg, term, env, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(cfg, term, env, curSt)
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
		cfg[finChnl.ID] = finChnl
		s.log.Debug("transition taking succeeded")
		spec.Term = wait.Cont
		return s.takeRecur(env, cfg, part, spec)
	case step.WaitSpec:
		curChnl, ok := cfg[term.X]
		if !ok {
			err = chnl.ErrDoesNotExist(term.X)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if curChnl.St == nil {
			err = chnl.ErrAlreadyClosed
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, ok := env[curChnl.StID]
		if !ok {
			err = state.ErrDoesNotExist(curChnl.StID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(cfg, term, env, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(cfg, term, env, curSt)
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
		_, ok = msg.Val.(step.CloseSpec)
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
		cfg[finChnl.ID] = finChnl
		spec.Term = term.Cont
		return s.takeRecur(env, cfg, part, spec)
	case step.SendSpec:
		curChnl, ok := cfg[term.A]
		if !ok {
			err = chnl.ErrDoesNotExist(term.A)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if curChnl.St == nil {
			err = chnl.ErrAlreadyClosed
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, ok := env[curChnl.StID]
		if !ok {
			err = state.ErrDoesNotExist(curChnl.StID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(cfg, term, env, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(cfg, term, env, curSt)
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
			StID:  curSt.(state.Prod).Next().RID(),
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
		cfg[newChnl.ID] = newChnl
		s.log.Debug("transition taking succeeded")
		recv.Cont = step.Subst(recv.Cont, recv.X, newChnl.ID)
		recv.Cont = step.Subst(recv.Cont, recv.Y, term.B)
		spec.Term = recv.Cont
		return s.takeRecur(env, cfg, part, spec)
	case step.RecvSpec:
		curChnl, ok := cfg[term.X]
		if !ok {
			err = chnl.ErrDoesNotExist(term.X)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if curChnl.St == nil {
			err = chnl.ErrAlreadyClosed
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("chnl", curChnl),
			)
			return err
		}
		curSt, ok := env[curChnl.StID]
		if !ok {
			err = state.ErrDoesNotExist(curChnl.StID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		if spec.AgentAK == part.PAK {
			err = s.checkProducer(cfg, term, env, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(cfg, term, env, curSt)
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
			StID:  curSt.(state.Prod).Next().RID(),
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
		cfg[newChnl.ID] = newChnl
		term.Cont = step.Subst(term.Cont, term.X, newChnl.ID)
		term.Cont = step.Subst(term.Cont, term.Y, val.B)
		spec.Term = term.Cont
		// TODO возможно здесь нужно начинать заново
		return s.takeRecur(env, cfg, part, spec)
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
			err = errAlreadyClosed(chnl.ConvertRootToRef(curChnl))
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
			err = s.checkProducer(cfg, term, env, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(cfg, term, env, curSt)
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
		cfg[newChnl.ID] = newChnl
		s.log.Debug("transition taking succeeded")
		spec.Term = step.Subst(cont.Conts[term.L], cont.Z, newChnl.ID)
		return s.takeRecur(env, cfg, part, spec)
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
			err = errAlreadyClosed(chnl.ConvertRootToRef(curChnl))
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
			err = s.checkProducer(cfg, term, env, curSt)
		} else if spec.AgentAK == part.CAK {
			err = s.checkConsumer(cfg, term, env, curSt)
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
		cfg[newChnl.ID] = newChnl
		s.log.Debug("transition taking succeeded")
		spec.Term = step.Subst(term.Conts[val.L], term.Z, newChnl.ID)
		return s.takeRecur(env, cfg, part, spec)
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
	// Consumable Access Key
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
func (s *dealService) checkProducer(
	cfg map[chnl.ID]chnl.Root,
	t step.Term,
	env map[state.ID]state.Root,
	st state.Root,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		return checkProducer(st, state.OneRoot{})
	case step.WaitSpec:
		return checkProducer(st, state.OneRoot{})
	case step.SendSpec:
		want, ok := st.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.TensorRoot{}, st)
		}
		// check value
		gotCh, ok := cfg[term.B]
		if !ok {
			return chnl.ErrDoesNotExist(term.B)
		}
		gotSt, ok := env[gotCh.StID]
		if !ok {
			return state.ErrDoesNotExist(gotCh.StID)
		}
		err := checkProducer(gotSt, want.B)
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
		gotCh, ok := cfg[term.Y]
		if !ok {
			return chnl.ErrDoesNotExist(term.Y)
		}
		gotSt, ok := env[gotCh.StID]
		if !ok {
			return state.ErrDoesNotExist(gotCh.StID)
		}
		err := checkProducer(gotSt, want.Y)
		if err != nil {
			return err
		}
		// check cont
		return s.checkProducer(cfg, term.Cont, env, want.Z)
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
		for wantLab, wantChoice := range want.Choices {
			gotBr, ok := term.Conts[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := s.checkProducer(cfg, gotBr, env, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(step.ErrUnexpectedTerm(t))
	}
}

func (s *dealService) checkConsumer(
	cfg map[chnl.ID]chnl.Root,
	t step.Term,
	env map[state.ID]state.Root,
	st state.Root,
) error {
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
		gotCh, ok := cfg[term.B]
		if !ok {
			return chnl.ErrDoesNotExist(term.B)
		}
		gotSt, ok := env[gotCh.StID]
		if !ok {
			return state.ErrDoesNotExist(gotCh.StID)
		}
		err := checkProducer(gotSt, want.Y)
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
		gotCh, ok := cfg[term.Y]
		if !ok {
			return chnl.ErrDoesNotExist(term.Y)
		}
		gotSt, ok := env[gotCh.StID]
		if !ok {
			return state.ErrDoesNotExist(gotCh.StID)
		}
		err := checkProducer(gotSt, want.B)
		if err != nil {
			return err
		}
		// check cont
		return s.checkConsumer(cfg, term.Cont, env, want.C)
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
			return fmt.Errorf("state mismatch: want %v choices, got %v conts", len(want.Choices), len(term.Conts))
		}
		for wantL, wantChoice := range want.Choices {
			gotCont, ok := term.Conts[wantL]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantL)
			}
			err := s.checkConsumer(cfg, gotCont, env, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(step.ErrUnexpectedTerm(t))
	}
}

func checkState(got, want state.Ref) error {
	if got != want {
		return fmt.Errorf("state mismatch: want %+v, got %+v", want, got)
	}
	return nil
}

// aka eqtp
func checkProducer(got, want state.Root) error {
	switch wantSt := want.(type) {
	case state.OneRoot:
		_, ok := got.(state.OneRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		return nil
	case state.TensorRoot:
		gotSt, ok := got.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		err := checkState(gotSt.B, wantSt.B)
		if err != nil {
			return err
		}
		return checkProducer(gotSt.C, wantSt.C)
	case state.LolliRoot:
		gotSt, ok := got.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		err := checkState(gotSt.Y, wantSt.Y)
		if err != nil {
			return err
		}
		return checkProducer(gotSt.Z, wantSt.Z)
	case state.PlusRoot:
		gotSt, ok := got.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v choices", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantLab, wantChoice := range wantSt.Choices {
			gotChoice, ok := gotSt.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := checkProducer(gotChoice, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	case state.WithRoot:
		gotSt, ok := got.(state.WithRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v choices", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantLab, wantChoice := range wantSt.Choices {
			gotChoice, ok := gotSt.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := checkProducer(gotChoice, wantChoice)
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
	case state.OneRoot:
		_, ok := got.(state.OneRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		return nil
	case state.TensorRoot:
		gotSt, ok := got.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		err := checkConsumer(gotSt.Y, wantSt.B)
		if err != nil {
			return err
		}
		return checkConsumer(gotSt.Z, wantSt.C)
	case state.LolliRoot:
		gotSt, ok := got.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		err := checkConsumer(gotSt.B, wantSt.Y)
		if err != nil {
			return err
		}
		return checkConsumer(gotSt.C, wantSt.Z)
	case state.PlusRoot:
		gotSt, ok := got.(state.WithRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v choices", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantLab, wantChoice := range wantSt.Choices {
			gotChoice, ok := gotSt.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := checkConsumer(gotChoice, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	case state.WithRoot:
		gotSt, ok := got.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v choices", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantLab, wantChoice := range wantSt.Choices {
			gotChoice, ok := gotSt.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := checkConsumer(gotChoice, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(state.ErrUnexpectedRoot(want))
	}
}

func errAlreadyClosed(ref chnl.Ref) error {
	return fmt.Errorf("channel already closed %+v", ref)
}
