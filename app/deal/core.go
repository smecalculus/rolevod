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
	Involve(PartSpec) (step.ProcRoot, error)
	Take(TranSpec) error
}

type dealService struct {
	deals    dealRepo
	seats    seat.SeatApi
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
	chnls chnl.Repo,
	procs step.Repo[step.ProcRoot],
	msgs step.Repo[step.MsgRoot],
	srvs step.Repo[step.SrvRoot],
	states state.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *dealService {
	name := slog.String("name", "dealService")
	return &dealService{
		deals, seats, chnls, procs, msgs, srvs, states, kinships, l.With(name),
	}
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

func (s *dealService) Involve(spec PartSpec) (step.ProcRoot, error) {
	s.log.Debug("seat involvement started", slog.Any("spec", spec))
	seat, err := s.seats.Retrieve(spec.SeatID)
	if err != nil {
		s.log.Error("seat selection failed",
			slog.Any("reason", err),
			slog.Any("spec", spec),
		)
		return step.ProcRoot{}, err
	}
	if len(seat.Ctx) != len(spec.Ctx) {
		err = fmt.Errorf("ctx mismatch: want %v items, got %v items", len(seat.Ctx), len(spec.Ctx))
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
			slog.Any("spec", spec),
		)
		return step.ProcRoot{}, err
	}
	newCtx := make(map[chnl.Name]chnl.ID, len(spec.Ctx))
	if len(spec.Ctx) > 0 {
		curCtx, err := s.chnls.SelectMany(maps.Values(spec.Ctx))
		if err != nil {
			s.log.Error("ctx selection failed",
				slog.Any("reason", err),
				slog.Any("ctx", spec.Ctx),
			)
			return step.ProcRoot{}, err
		}
		for i, got := range curCtx {
			// TODO обеспечить порядок
			// TODO проверять по значению, а не по ссылке
			err = checkState(got.St, seat.Ctx[i].St)
			if err != nil {
				s.log.Error("type checking failed",
					slog.Any("reason", err),
					slog.Any("spec", spec),
				)
				return step.ProcRoot{}, err
			}
		}
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
			return step.ProcRoot{}, err
		}
		for _, ch := range chnls {
			newCtx[ch.Name] = ch.ID
		}
	}
	via := chnl.Root{
		ID:   id.New(),
		Name: seat.Via.Name,
		StID: seat.Via.StID,
		St:   seat.Via.St,
	}
	err = s.chnls.Insert(via)
	if err != nil {
		s.log.Error("via insertion failed",
			slog.Any("reason", err),
			slog.Any("via", via),
		)
		return step.ProcRoot{}, err
	}
	proc := step.ProcRoot{
		ID:  id.New(),
		PID: via.ID,
		Ctx: newCtx,
		Term: step.CTASpec{
			Key:  ak.New(),
			Seat: spec.SeatID,
		},
	}
	err = s.procs.Insert(proc)
	if err != nil {
		s.log.Error("process insertion failed",
			slog.Any("reason", err),
			slog.Any("proc", proc),
		)
		return step.ProcRoot{}, err
	}
	s.log.Debug("seat involvement succeeded", slog.Any("proc", proc))
	return proc, nil
}

func (s *dealService) Take(spec TranSpec) error {
	if spec.Term == nil {
		panic(step.ErrUnexpectedTermType(spec.Term))
	}
	s.log.Debug("transition taking started", slog.Any("spec", spec))
	// proc checking
	proc, err := s.procs.SelectByPID(spec.PID)
	if err != nil {
		s.log.Error("process selection failed",
			slog.Any("reason", err),
			slog.Any("id", spec.PID),
		)
		return err
	}
	if proc == nil {
		err = step.ErrDoesNotExist(spec.PID)
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
		)
		return err
	}
	_, ok := proc.Term.(step.CTASpec)
	if !ok {
		err = step.ErrTermMismatch(spec.Term, step.CTASpec{})
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
		)
		return err
	}
	// if spec.Key != cta.Key {
	// 	err = ak.ErrUnexpectedKey(spec.Key)
	// 	s.log.Error("transition taking failed",
	// 		slog.Any("reason", err),
	// 	)
	// 	return err
	// }
	// type checking
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
	if spec.Term.Via() == spec.PID {
		s.checkProvider(env, cfg, spec.Term)
	} else {
		s.checkClient(env, cfg, spec.Term)
	}
	// step taking
	proc.Term = spec.Term
	return s.takeProcWith(*proc, cfg, env)
}

func (s *dealService) takeProc(
	proc step.ProcRoot,
) (err error) {
	s.log.Debug("transition taking started", slog.Any("proc", proc))
	viaID := proc.Term.Via().(chnl.ID)
	cfg, err := s.chnls.SelectCfg([]chnl.ID{viaID})
	if err != nil {
		s.log.Error("channel selection failed",
			slog.Any("reason", err),
			slog.Any("id", viaID),
		)
		return err
	}
	stID := cfg[viaID].StID
	env, err := s.states.SelectEnv([]state.ID{stID})
	if err != nil {
		s.log.Error("state selection failed",
			slog.Any("reason", err),
			slog.Any("id", stID),
		)
		return err
	}
	return s.takeProcWith(proc, cfg, env)
}

func (s *dealService) takeProcWith(
	proc step.ProcRoot,
	cfg map[chnl.ID]chnl.Root,
	env map[state.ID]state.Root,
) (err error) {
	switch term := proc.Term.(type) {
	case step.CloseSpec:
		viaID, ok := term.A.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.A)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg[viaID]
		if !ok {
			err := chnl.ErrDoesNotExist(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		srv, err := s.srvs.SelectByVID(viaID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if srv == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				PID: proc.PID,
				VID: curVia.ID,
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
		// consume and close channel
		finVia := chnl.Root{
			ID:    id.New(),
			Name:  curVia.Name,
			PreID: curVia.ID,
			St:    nil,
		}
		err = s.chnls.Insert(finVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", finVia),
			)
			return err
		}
		newProc := step.ProcRoot{
			ID:   id.New(),
			PID:  srv.PID,
			Term: wait.Cont,
		}
		s.log.Debug("transition taking succeeded")
		return s.takeProc(newProc)
	case step.WaitSpec:
		viaID, ok := term.X.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.X)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg[viaID]
		if !ok {
			err = chnl.ErrDoesNotExist(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		msg, err := s.msgs.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if msg == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				PID:  proc.PID,
				VID:  curVia.ID,
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
		if ok {
			// consume and close channel
			finVia := chnl.Root{
				ID:    id.New(),
				Name:  curVia.Name,
				PreID: curVia.ID,
				St:    nil,
			}
			err = s.chnls.Insert(finVia)
			if err != nil {
				s.log.Error("channel insertion failed",
					slog.Any("reason", err),
					slog.Any("chnl", finVia),
				)
				return err
			}
			newProc := step.ProcRoot{
				ID:   id.New(),
				PID:  proc.PID,
				Term: term.Cont,
			}
			s.log.Debug("transition taking succeeded")
			return s.takeProc(newProc)
		}
		fwd, ok := msg.Val.(step.FwdSpec)
		if ok {
			newTerm := step.Subst(term, fwd.C, fwd.D.(chnl.ID))
			newProc := step.ProcRoot{
				ID:   id.New(),
				PID:  proc.PID,
				Term: newTerm,
			}
			s.log.Debug("transition taking succeeded")
			return s.takeProc(newProc)
		}
		err = fmt.Errorf("unexpected val type: want %T, got %T", step.CloseSpec{}, msg.Val)
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
			slog.Any("val", msg.Val),
		)
		return err
	case step.SendSpec:
		viaID, ok := term.A.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.A)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg[viaID]
		if !ok {
			err = chnl.ErrDoesNotExist(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		srv, err := s.srvs.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if srv == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				PID: proc.PID,
				VID: curVia.ID,
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
		st := env[curVia.StID]
		newVia := chnl.Root{
			ID:    id.New(),
			Name:  curVia.Name,
			PreID: curVia.ID,
			StID:  st.(state.Prod).Next().RID(),
			St:    st.(state.Prod).Next(),
		}
		err = s.chnls.Insert(newVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", newVia),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		recv.Cont = step.Subst(recv.Cont, recv.X, newVia.ID)
		recv.Cont = step.Subst(recv.Cont, recv.Y, term.B.(chnl.ID))
		newProc := step.ProcRoot{
			ID:   id.New(),
			PID:  srv.PID,
			Term: recv.Cont,
		}
		return s.takeProc(newProc)
	case step.RecvSpec:
		viaID, ok := term.X.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.X)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg[viaID]
		if !ok {
			err = chnl.ErrDoesNotExist(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		msg, err := s.msgs.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if msg == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				PID:  proc.PID,
				VID:  curVia.ID,
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
		curSt := env[curVia.StID]
		newVia := chnl.Root{
			ID:    id.New(),
			Name:  curVia.Name,
			PreID: curVia.ID,
			StID:  curSt.(state.Prod).Next().RID(),
			St:    curSt.(state.Prod).Next(),
		}
		err = s.chnls.Insert(newVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", newVia),
			)
			return err
		}
		term.Cont = step.Subst(term.Cont, term.X, newVia.ID)
		term.Cont = step.Subst(term.Cont, term.Y, val.B.(chnl.ID))
		newProc := step.ProcRoot{
			ID:   id.New(),
			PID:  proc.PID,
			Term: term.Cont,
		}
		return s.takeProc(newProc)
	case step.LabSpec:
		viaID, ok := term.C.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.C)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg[viaID]
		if !ok {
			err = chnl.ErrDoesNotExist(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		srv, err := s.srvs.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if srv == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				PID: proc.PID,
				VID: curVia.ID,
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
		curSt := env[curVia.StID]
		newVia := chnl.Root{
			ID:    id.New(),
			Name:  curVia.Name,
			PreID: curVia.ID,
			StID:  curSt.(state.Sum).Next(term.L).RID(),
			St:    curSt.(state.Sum).Next(term.L),
		}
		err = s.chnls.Insert(newVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", newVia),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		newProc := step.ProcRoot{
			ID:   id.New(),
			PID:  srv.PID,
			Term: step.Subst(cont.Conts[term.L], cont.Z, newVia.ID),
		}
		return s.takeProc(newProc)
	case step.CaseSpec:
		viaID, ok := term.Z.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.Z)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg[viaID]
		if !ok {
			err = chnl.ErrDoesNotExist(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		msg, err := s.msgs.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if msg == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				PID:  proc.PID,
				VID:  curVia.ID,
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
		curSt := env[curVia.StID]
		newVia := chnl.Root{
			ID:    id.New(),
			Name:  curVia.Name,
			PreID: curVia.ID,
			StID:  curSt.(state.Sum).Next(val.L).RID(),
			St:    curSt.(state.Sum).Next(val.L),
		}
		err = s.chnls.Insert(newVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", newVia),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		newProc := step.ProcRoot{
			ID:   id.New(),
			PID:  proc.PID,
			Term: step.Subst(term.Conts[val.L], term.Z, newVia.ID),
		}
		return s.takeProc(newProc)
	case step.SpawnSpec:
		newProc, err := s.Involve(PartSpec{SeatID: term.SeatID, Ctx: term.Ctx})
		if err != nil {
			return err
		}
		term.Cont = step.Subst(term.Cont, term.Z, newProc.PID)
		return s.Take(TranSpec{PID: proc.PID, Term: term.Cont})
	case step.FwdSpec:
		viaID, ok := term.C.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.C)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg[viaID]
		if !ok {
			err = chnl.ErrDoesNotExist(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSt, ok := env[curVia.StID]
		if !ok {
			err = state.ErrDoesNotExist(curVia.StID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		switch curSt.Pol() {
		case state.Pos:
			srv, err := s.srvs.SelectByVID(curVia.ID)
			if err != nil {
				s.log.Error("service selection failed",
					slog.Any("reason", err),
					slog.Any("via", curVia),
				)
				return err
			}
			if srv == nil {
				newMsg := step.MsgRoot{
					ID:  id.New(),
					PID: proc.PID,
					VID: curVia.ID,
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
			newCont := step.Subst(srv.Cont, term.C, term.D.(chnl.ID))
			newProc := step.ProcRoot{
				ID:   id.New(),
				PID:  srv.PID,
				Term: newCont,
			}
			return s.takeProc(newProc)
		case state.Neg:
			msg, err := s.msgs.SelectByVID(curVia.ID)
			if err != nil {
				s.log.Error("message selection failed",
					slog.Any("reason", err),
					slog.Any("via", curVia),
				)
				return err
			}
			if msg == nil {
				newSrv := step.SrvRoot{
					ID:   id.New(),
					PID:  proc.PID,
					VID:  curVia.ID,
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
			newVal := step.Subst(msg.Val, term.C, term.D.(chnl.ID))
			newProc := step.ProcRoot{
				ID:   id.New(),
				PID:  msg.PID,
				Term: newVal,
			}
			return s.takeProc(newProc)
		default:
			panic(state.ErrUnexpectedPolarity(curSt.Pol()))
		}
		return nil
	default:
		panic(step.ErrUnexpectedTermType(proc.Term))
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

// Participation aka lightweight Spawn
type PartSpec struct {
	DealID ID
	SeatID seat.ID
	Ctx    map[chnl.Name]chnl.ID
}

// Transition
type TranSpec struct {
	// Deal ID
	DID ID
	// Proc ID
	PID chnl.ID
	// Agent Access Key
	Key  ak.ADT
	Term step.Term
}

// aka checkExp
func (s *dealService) checkProvider(
	env map[state.ID]state.Root,
	cfg map[chnl.ID]chnl.Root,
	t step.Term,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		// check via
		gotA, err := findState(env, cfg, term.A.(chnl.ID))
		if err != nil {
			return err
		}
		return checkProvider(gotA, state.OneRoot{})
	case step.WaitSpec:
		// check via
		gotX, err := findState(env, cfg, term.X.(chnl.ID))
		if err != nil {
			return err
		}
		_, ok := gotX.(state.OneRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.OneRoot{}, gotX)
		}
		// check cont
		return s.checkProvider(env, cfg, term.Cont)
	case step.SendSpec:
		// check via
		gotA, err := findState(env, cfg, term.A.(chnl.ID))
		if err != nil {
			return err
		}
		want, ok := gotA.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.TensorRoot{}, gotA)
		}
		// check value
		gotB, err := findState(env, cfg, term.B.(chnl.ID))
		if err != nil {
			s.log.Error("type checking failed",
				slog.Any("reason", err),
			)
			return err
		}
		err = checkProvider(gotB, want.B)
		if err != nil {
			return err
		}
		// no cont to check
		return nil
	case step.RecvSpec:
		// check via
		gotX, err := findState(env, cfg, term.X.(chnl.ID))
		if err != nil {
			return err
		}
		want, ok := gotX.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.LolliRoot{}, gotX)
		}
		// check value
		gotY, err := findState(env, cfg, term.Y.(chnl.ID))
		if err != nil {
			return err
		}
		err = checkProvider(gotY, want.Y)
		if err != nil {
			return err
		}
		// check cont
		return s.checkProvider(env, cfg, term.Cont)
	case step.LabSpec:
		// check via
		gotC, err := findState(env, cfg, term.C.(chnl.ID))
		if err != nil {
			return err
		}
		want, ok := gotC.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.PlusRoot{}, gotC)
		}
		// check label
		_, ok = want.Choices[term.L]
		if !ok {
			return fmt.Errorf("label mismatch: want %q, got nothing", term.L)
		}
		// no cont to check
		return nil
	case step.CaseSpec:
		// check via
		gotZ, err := findState(env, cfg, term.Z.(chnl.ID))
		if err != nil {
			return err
		}
		want, ok := gotZ.(state.WithRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.WithRoot{}, gotZ)
		}
		// check conts
		if len(term.Conts) != len(want.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v conts", len(want.Choices), len(term.Conts))
		}
		for wantLab := range want.Choices {
			gotCont, ok := term.Conts[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := s.checkProvider(env, cfg, gotCont)
			if err != nil {
				return err
			}
		}
		return nil
	case step.SpawnSpec:
		return nil
	case step.FwdSpec:
		return nil
	default:
		panic(step.ErrUnexpectedTermType(t))
	}
}

func (s *dealService) checkClient(
	env map[state.ID]state.Root,
	cfg map[chnl.ID]chnl.Root,
	t step.Term,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		// check via
		gotA, err := findState(env, cfg, term.A.(chnl.ID))
		if err != nil {
			return err
		}
		return checkClient(gotA, state.OneRoot{})
	case step.WaitSpec:
		// check via
		gotX, err := findState(env, cfg, term.X.(chnl.ID))
		if err != nil {
			return err
		}
		_, ok := gotX.(state.OneRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.OneRoot{}, gotX)
		}
		// check cont
		return s.checkClient(env, cfg, term.Cont)
	case step.SendSpec:
		// check via
		gotA, err := findState(env, cfg, term.A.(chnl.ID))
		if err != nil {
			return err
		}
		want, ok := gotA.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.LolliRoot{}, gotA)
		}
		// check value
		gotB, err := findState(env, cfg, term.B.(chnl.ID))
		if err != nil {
			s.log.Error("type checking failed",
				slog.Any("reason", err),
			)
			return err
		}
		err = checkProvider(gotB, want.Y)
		if err != nil {
			return err
		}
		// no cont to check
		return nil
	case step.RecvSpec:
		// check via
		gotX, err := findState(env, cfg, term.X.(chnl.ID))
		if err != nil {
			return err
		}
		want, ok := gotX.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.TensorRoot{}, gotX)
		}
		// check value
		gotY, err := findState(env, cfg, term.Y.(chnl.ID))
		if err != nil {
			return err
		}
		err = checkProvider(gotY, want.B)
		if err != nil {
			return err
		}
		// check cont
		return s.checkClient(env, cfg, term.Cont)
	case step.LabSpec:
		// check via
		gotC, err := findState(env, cfg, term.C.(chnl.ID))
		if err != nil {
			return err
		}
		want, ok := gotC.(state.WithRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.WithRoot{}, gotC)
		}
		// check label
		_, ok = want.Choices[term.L]
		if !ok {
			return fmt.Errorf("label mismatch: want %q, got nothing", term.L)
		}
		// no cont to check
		return nil
	case step.CaseSpec:
		// check via
		gotZ, err := findState(env, cfg, term.Z.(chnl.ID))
		if err != nil {
			return err
		}
		want, ok := gotZ.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", state.PlusRoot{}, gotZ)
		}
		// check conts
		if len(term.Conts) != len(want.Choices) {
			return fmt.Errorf("state mismatch: want %v choices, got %v conts", len(want.Choices), len(term.Conts))
		}
		for wantLab := range want.Choices {
			gotCont, ok := term.Conts[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := s.checkClient(env, cfg, gotCont)
			if err != nil {
				return err
			}
		}
		return nil
	case step.SpawnSpec:
		return nil
	case step.FwdSpec:
		return nil
	default:
		panic(step.ErrUnexpectedTermType(t))
	}
}

func checkState(got, want state.Ref) error {
	if got != want {
		return fmt.Errorf("state mismatch: want %+v, got %+v", want, got)
	}
	return nil
}

// aka eqtp
func checkProvider(got, want state.Root) error {
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
		err := checkProvider(gotSt.B, wantSt.B)
		if err != nil {
			return err
		}
		return checkProvider(gotSt.C, wantSt.C)
	case state.LolliRoot:
		gotSt, ok := got.(state.LolliRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		err := checkProvider(gotSt.Y, wantSt.Y)
		if err != nil {
			return err
		}
		return checkProvider(gotSt.Z, wantSt.Z)
	case state.PlusRoot:
		gotSt, ok := got.(state.PlusRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("choices mismatch: want %v, got %v", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantLab, wantChoice := range wantSt.Choices {
			gotChoice, ok := gotSt.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := checkProvider(gotChoice, wantChoice)
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
			return fmt.Errorf("choices mismatch: want %v, got %v", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantLab, wantChoice := range wantSt.Choices {
			gotChoice, ok := gotSt.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := checkProvider(gotChoice, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(state.ErrUnexpectedRootType(want))
	}
}

func checkClient(got, want state.Root) error {
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
		err := checkClient(gotSt.Y, wantSt.B)
		if err != nil {
			return err
		}
		return checkClient(gotSt.Z, wantSt.C)
	case state.LolliRoot:
		gotSt, ok := got.(state.TensorRoot)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %T", want, got)
		}
		err := checkClient(gotSt.B, wantSt.Y)
		if err != nil {
			return err
		}
		return checkClient(gotSt.C, wantSt.Z)
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
			err := checkClient(gotChoice, wantChoice)
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
			err := checkClient(gotChoice, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(state.ErrUnexpectedRootType(want))
	}
}

func findState(
	env map[state.ID]state.Root,
	cfg map[chnl.ID]chnl.Root,
	chID chnl.ID,
) (state.Root, error) {
	gotCh, ok := cfg[chID]
	if !ok {
		return nil, chnl.ErrDoesNotExist(chID)
	}
	if gotCh.St == nil {
		return nil, chnl.ErrAlreadyClosed(chID)
	}
	gotSt, ok := env[gotCh.StID]
	if !ok {
		return nil, state.ErrDoesNotExist(gotCh.StID)
	}
	return gotSt, nil
}
