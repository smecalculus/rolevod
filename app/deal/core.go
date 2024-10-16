package deal

import (
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
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
	steps    step.Repo
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newDealService(
	deals dealRepo,
	seats seat.SeatApi,
	chnls chnl.Repo,
	steps step.Repo,
	states state.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *dealService {
	name := slog.String("name", "dealService")
	return &dealService{
		deals, seats, chnls, steps, states, kinships, l.With(name),
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

func (s *dealService) Involve(gotSpec PartSpec) (step.ProcRoot, error) {
	s.log.Debug("seat involvement started", slog.Any("spec", gotSpec))
	wantSpec, err := s.seats.Retrieve(gotSpec.SeatID)
	if err != nil {
		s.log.Error("seat selection failed",
			slog.Any("reason", err),
			slog.Any("id", gotSpec.SeatID),
		)
		return step.ProcRoot{}, err
	}
	if len(gotSpec.Ctx) != len(wantSpec.Ctx) {
		err := fmt.Errorf("context mismatch: want %v items, got %v items", len(wantSpec.Ctx), len(gotSpec.Ctx))
		s.log.Error("seat involvement failed",
			slog.Any("reason", err),
			slog.Any("ctx", gotSpec.Ctx),
		)
		return step.ProcRoot{}, err
	}
	newVia := chnl.Root{
		ID:   id.New(),
		Name: wantSpec.Via.Name,
		StID: wantSpec.Via.StID,
	}
	err = s.chnls.Insert(newVia)
	if err != nil {
		s.log.Error("via insertion failed",
			slog.Any("reason", err),
			slog.Any("via", newVia),
		)
		return step.ProcRoot{}, err
	}
	if len(gotSpec.Ctx) > 0 {
		gotCtx, err := s.chnls.SelectCtx(gotSpec.OID, gotSpec.Ctx)
		if err != nil {
			s.log.Error("context selection failed",
				slog.Any("reason", err),
				slog.Any("ctx", gotSpec.Ctx),
			)
			return step.ProcRoot{}, err
		}
		for i, got := range gotCtx {
			// TODO обеспечить порядок
			// TODO проверять по значению, а не по ссылке
			err = state.CheckRef(got.StID, wantSpec.Ctx[i].StID)
			if err != nil {
				s.log.Error("type checking failed",
					slog.Any("reason", err),
				)
				return step.ProcRoot{}, err
			}
		}
		err = s.chnls.Transfer(gotSpec.OID, newVia.ID, gotSpec.Ctx)
		if err != nil {
			s.log.Error("context transfer failed",
				slog.Any("reason", err),
				slog.Any("from", gotSpec.OID),
				slog.Any("to", newVia.ID),
				slog.Any("ctx", gotSpec.Ctx),
			)
			return step.ProcRoot{}, err
		}
	}
	newProc := step.ProcRoot{
		ID:  id.New(),
		PID: newVia.ID,
		Ctx: gotSpec.Ctx,
		Term: step.CTASpec{
			AK:  ak.New(),
			SID: gotSpec.SeatID,
		},
	}
	err = s.steps.Insert(newProc)
	if err != nil {
		s.log.Error("process insertion failed",
			slog.Any("reason", err),
			slog.Any("proc", newProc),
		)
		return step.ProcRoot{}, err
	}
	s.log.Debug("seat involvement succeeded", slog.Any("proc", newProc))
	return newProc, nil
}

func (s *dealService) Take(spec TranSpec) error {
	if spec.Term == nil {
		panic(step.ErrTermValueNil(spec.PID))
	}
	s.log.Debug("transition taking started", slog.Any("spec", spec))
	// proc checking
	curSem, err := s.steps.SelectByPID(spec.PID)
	if err != nil {
		s.log.Error("process selection failed",
			slog.Any("reason", err),
			slog.Any("id", spec.PID),
		)
		return err
	}
	if curSem == nil {
		err = step.ErrDoesNotExist(spec.PID)
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
		)
		return err
	}
	proc, ok := curSem.(step.ProcRoot)
	if !ok {
		err = step.ErrRootTypeMismatch(curSem, step.ProcRoot{})
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
		)
		return err
	}
	_, ok = proc.Term.(step.CTASpec)
	if !ok {
		err = step.ErrTermTypeMismatch(spec.Term, step.CTASpec{})
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
		)
		return err
	}
	providable, err := s.chnls.SelectByID(proc.PID)
	if err != nil {
		s.log.Error("zc selection failed",
			slog.Any("reason", err),
			slog.Any("id", proc.PID),
		)
		return err
	}
	consumables, err := s.chnls.SelectCtx(proc.PID, proc.Ctx)
	if err != nil {
		s.log.Error("context selection failed",
			slog.Any("reason", err),
			slog.Any("pid", proc.PID),
			slog.Any("ctx", proc.Ctx),
		)
		return err
	}
	// type checking
	stIDs := chnl.CollectStIDs(append(consumables, providable))
	states, err := s.states.SelectEnv(stIDs)
	if err != nil {
		s.log.Error("env selection failed",
			slog.Any("reason", err),
			slog.Any("ids", stIDs),
		)
		return err
	}
	ctx := join(consumables, states)
	zc := state.ZC{Z: providable.ID, C: states[providable.StID]}
	if spec.PID == spec.Term.Via() {
		s.checkProvider(ctx, zc, spec.Term)
	} else {
		s.checkClient(ctx, zc, spec.Term)
	}
	// step taking
	cfg := convertToCfg(append(consumables, providable))
	proc.Term = spec.Term
	return s.takeProcWith(proc, cfg, states)
}

func (s *dealService) takeProc(
	proc step.ProcRoot,
) (err error) {
	s.log.Debug("transition taking started", slog.Any("proc", proc))
	zc, err := s.chnls.SelectByID(proc.PID)
	if err != nil {
		s.log.Error("zc selection failed",
			slog.Any("reason", err),
			slog.Any("pid", proc.PID),
		)
		return err
	}
	ctx, err := s.chnls.SelectCtx(proc.PID, proc.Ctx)
	if err != nil {
		s.log.Error("context selection failed",
			slog.Any("reason", err),
			slog.Any("ctx", proc.Ctx),
		)
		return err
	}
	stIDs := chnl.CollectStIDs(append(ctx, zc))
	env, err := s.states.SelectEnv(stIDs)
	if err != nil {
		s.log.Error("env selection failed",
			slog.Any("reason", err),
			slog.Any("ids", stIDs),
		)
		return err
	}
	cfg := convertToCfg(append(ctx, zc))
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
			err := chnl.ErrMissingInCfg(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSem, err := s.steps.SelectByVID(viaID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if curSem == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				PID: proc.PID,
				VID: curVia.ID,
				Val: term,
			}
			err := s.steps.Insert(newMsg)
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
		srv, ok := curSem.(step.SrvRoot)
		if !ok {
			err = step.ErrRootTypeUnexpected(curSem)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		wait, ok := srv.Cont.(step.WaitSpec)
		if !ok {
			err = fmt.Errorf("cont type mismatch: want %T, got %T", step.WaitSpec{}, srv.Cont)
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
			StID:  id.Empty(),
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
			err = chnl.ErrMissingInCfg(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSem, err := s.steps.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if curSem == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				PID:  proc.PID,
				VID:  curVia.ID,
				Cont: term,
			}
			err = s.steps.Insert(newSrv)
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
		msg, ok := curSem.(step.MsgRoot)
		if !ok {
			err = step.ErrRootTypeMismatch(curSem, step.MsgRoot{})
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		var newProc step.ProcRoot
		switch val := msg.Val.(type) {
		case step.CloseSpec:
			// consume and close channel
			finVia := chnl.Root{
				ID:    id.New(),
				Name:  curVia.Name,
				PreID: curVia.ID,
				StID:  id.Empty(),
			}
			err = s.chnls.Insert(finVia)
			if err != nil {
				s.log.Error("channel insertion failed",
					slog.Any("reason", err),
					slog.Any("via", finVia),
				)
				return err
			}
			newProc = step.ProcRoot{
				ID:   id.New(),
				PID:  proc.PID,
				Term: term.Cont,
			}
		case step.FwdSpec:
			d, ok := val.D.(chnl.ID)
			if !ok {
				err = chnl.ErrNotAnID(val.D)
				s.log.Error("transition taking failed",
					slog.Any("reason", err),
				)
				return err
			}
			err := s.chnls.Transfer(msg.PID, proc.PID, []chnl.ID{d})
			if err != nil {
				s.log.Error("channel transfer failed",
					slog.Any("reason", err),
					slog.Any("id", d),
				)
				return err
			}
			newProc = step.ProcRoot{
				ID:   id.New(),
				PID:  proc.PID,
				Ctx:  []chnl.ID{d},
				Term: step.Subst(term, val.C, d),
			}
		default:
			panic(step.ErrValTypeUnexpected(msg.Val))
		}
		s.log.Debug("transition taking succeeded")
		return s.takeProc(newProc)
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
			err = chnl.ErrMissingInCfg(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSem, err := s.steps.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if curSem == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				PID: proc.PID,
				VID: curVia.ID,
				Val: term,
			}
			err = s.steps.Insert(newMsg)
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
		srv, ok := curSem.(step.SrvRoot)
		if !ok {
			err = step.ErrRootTypeMismatch(curSem, step.SrvRoot{})
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
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
		}
		err = s.chnls.Insert(newVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", newVia),
			)
			return err
		}
		b, ok := cfg[term.B.(chnl.ID)]
		if !ok {
			err = chnl.ErrMissingInCfg(b.ID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		err = s.chnls.Transfer(proc.PID, srv.PID, []chnl.ID{b.ID})
		if err != nil {
			s.log.Error("channel transfer failed",
				slog.Any("reason", err),
				slog.Any("id", b.ID),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		recv.Cont = step.Subst(recv.Cont, recv.X, newVia.ID)
		recv.Cont = step.Subst(recv.Cont, recv.Y, b.ID)
		newProc := step.ProcRoot{
			ID:   id.New(),
			PID:  chnl.Subst(srv.PID, curVia.ID, newVia.ID),
			Ctx:  []chnl.ID{b.ID},
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
			err = chnl.ErrMissingInCfg(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSem, err := s.steps.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if curSem == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				PID:  proc.PID,
				VID:  curVia.ID,
				Cont: term,
			}
			err = s.steps.Insert(newSrv)
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
		msg, ok := curSem.(step.MsgRoot)
		if !ok {
			err = step.ErrRootTypeMismatch(curSem, step.MsgRoot{})
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
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
		}
		err = s.chnls.Insert(newVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", newVia),
			)
			return err
		}
		b, ok := cfg[val.B.(chnl.ID)]
		if !ok {
			err = chnl.ErrMissingInCfg(b.ID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		err = s.chnls.Transfer(msg.PID, proc.PID, []chnl.ID{b.ID})
		if err != nil {
			s.log.Error("channel transfer failed",
				slog.Any("reason", err),
				slog.Any("id", b.ID),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		term.Cont = step.Subst(term.Cont, term.X, newVia.ID)
		term.Cont = step.Subst(term.Cont, term.Y, b.ID)
		newProc := step.ProcRoot{
			ID:   id.New(),
			PID:  chnl.Subst(proc.PID, curVia.ID, newVia.ID),
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
			err = chnl.ErrMissingInCfg(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSem, err := s.steps.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if curSem == nil {
			newMsg := step.MsgRoot{
				ID:  id.New(),
				PID: proc.PID,
				VID: curVia.ID,
				Val: term,
			}
			err = s.steps.Insert(newMsg)
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
		srv, ok := curSem.(step.SrvRoot)
		if !ok {
			err = step.ErrRootTypeMismatch(curSem, step.SrvRoot{})
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
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
			PID:  chnl.Subst(srv.PID, curVia.ID, newVia.ID),
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
			err = chnl.ErrMissingInCfg(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSem, err := s.steps.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("via", curVia),
			)
			return err
		}
		if curSem == nil {
			newSrv := step.SrvRoot{
				ID:   id.New(),
				PID:  proc.PID,
				VID:  curVia.ID,
				Cont: term,
			}
			err = s.steps.Insert(newSrv)
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
		msg, ok := curSem.(step.MsgRoot)
		if !ok {
			err = step.ErrRootTypeMismatch(curSem, step.MsgRoot{})
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		val, ok := msg.Val.(step.LabSpec)
		if !ok {
			err = fmt.Errorf("val type mismatch: want %T, got %T", step.LabSpec{}, msg.Val)
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
			PID:  chnl.Subst(proc.PID, curVia.ID, newVia.ID),
			Term: step.Subst(term.Conts[val.L], term.Z, newVia.ID),
		}
		return s.takeProc(newProc)
	case step.SpawnSpec:
		newProc, err := s.Involve(PartSpec{SeatID: term.SeatID, OID: proc.PID, Ctx: term.Ctx})
		if err != nil {
			return err
		}
		err = s.chnls.Transfer(id.Empty(), proc.PID, []chnl.ID{newProc.PID})
		if err != nil {
			s.log.Error("channel transfer failed",
				slog.Any("reason", err),
				slog.Any("id", newProc.PID),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		// TODO актуализировать контекст
		proc.Term = step.Subst(term.Cont, term.Z, newProc.PID)
		return s.takeProc(proc)
	case step.FwdSpec:
		viaID, ok := term.C.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.C)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		viaCh, ok := cfg[viaID]
		if !ok {
			err = chnl.ErrMissingInCfg(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSt, ok := env[viaCh.StID]
		if !ok {
			err = state.ErrMissingInEnv(viaCh.StID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSem, err := s.steps.SelectByVID(viaCh.ID)
		if err != nil {
			s.log.Error("step selection failed",
				slog.Any("reason", err),
				slog.Any("via", viaCh),
			)
			return err
		}
		switch curSt.Pol() {
		case state.Pos:
			switch sem := curSem.(type) {
			case step.SrvRoot:
				c, ok := cfg[term.C.(chnl.ID)]
				if !ok {
					err = chnl.ErrMissingInCfg(c.ID)
					s.log.Error("transition taking failed",
						slog.Any("reason", err),
					)
					return err
				}
				err := s.chnls.Transfer(proc.PID, sem.PID, []chnl.ID{c.ID})
				if err != nil {
					s.log.Error("channel transfer failed",
						slog.Any("reason", err),
						slog.Any("id", c.ID),
					)
					return err
				}
				newProc := step.ProcRoot{
					ID:   id.New(),
					PID:  sem.PID,
					Ctx:  []chnl.ID{c.ID},
					Term: step.Subst(sem.Cont, term.D, c.ID),
				}
				s.log.Debug("transition taking succeeded")
				return s.takeProc(newProc)
			case step.MsgRoot:
				d, ok := cfg[term.D.(chnl.ID)]
				if !ok {
					err = chnl.ErrMissingInCfg(d.ID)
					s.log.Error("transition taking failed",
						slog.Any("reason", err),
					)
					return err
				}
				err := s.chnls.Transfer(proc.PID, sem.PID, []chnl.ID{d.ID})
				if err != nil {
					s.log.Error("channel transfer failed",
						slog.Any("reason", err),
						slog.Any("id", d.ID),
					)
					return err
				}
				newProc := step.ProcRoot{
					ID:   id.New(),
					PID:  sem.PID,
					Ctx:  []chnl.ID{d.ID},
					Term: step.Subst(sem.Val, term.C, d.ID),
				}
				s.log.Debug("transition taking succeeded")
				return s.takeProc(newProc)
			case nil:
				newMsg := step.MsgRoot{
					ID:  id.New(),
					PID: proc.PID,
					VID: viaCh.ID,
					Val: term,
				}
				err := s.steps.Insert(newMsg)
				if err != nil {
					s.log.Error("message insertion failed",
						slog.Any("reason", err),
						slog.Any("msg", newMsg),
					)
					return err
				}
				s.log.Debug("transition taking half done", slog.Any("msg", newMsg))
				return nil
			default:
				panic(step.ErrRootTypeUnexpected(curSem))
			}
		case state.Neg:
			switch sem := curSem.(type) {
			case step.SrvRoot:
				d, ok := cfg[term.D.(chnl.ID)]
				if !ok {
					err = chnl.ErrMissingInCfg(d.ID)
					s.log.Error("transition taking failed",
						slog.Any("reason", err),
					)
					return err
				}
				err := s.chnls.Transfer(proc.PID, sem.PID, []chnl.ID{d.ID})
				if err != nil {
					s.log.Error("channel transfer failed",
						slog.Any("reason", err),
						slog.Any("id", d.ID),
					)
					return err
				}
				newProc := step.ProcRoot{
					ID:   id.New(),
					PID:  sem.PID,
					Ctx:  []chnl.ID{d.ID},
					Term: step.Subst(sem.Cont, term.C, d.ID),
				}
				s.log.Debug("transition taking succeeded")
				return s.takeProc(newProc)
			case step.MsgRoot:
				c := cfg[term.C.(chnl.ID)]
				if !ok {
					err = chnl.ErrMissingInCfg(c.ID)
					s.log.Error("transition taking failed",
						slog.Any("reason", err),
					)
					return err
				}
				err := s.chnls.Transfer(proc.PID, sem.PID, []chnl.ID{c.ID})
				if err != nil {
					s.log.Error("channel transfer failed",
						slog.Any("reason", err),
						slog.Any("id", c.ID),
					)
					return err
				}
				newProc := step.ProcRoot{
					ID:   id.New(),
					PID:  sem.PID,
					Ctx:  []chnl.ID{c.ID},
					Term: step.Subst(sem.Val, term.D, c.ID),
				}
				s.log.Debug("transition taking succeeded")
				return s.takeProc(newProc)
			case nil:
				newSrv := step.SrvRoot{
					ID:   id.New(),
					PID:  proc.PID,
					VID:  viaCh.ID,
					Cont: term,
				}
				err = s.steps.Insert(newSrv)
				if err != nil {
					s.log.Error("service insertion failed",
						slog.Any("reason", err),
						slog.Any("srv", newSrv),
					)
					return err
				}
				s.log.Debug("transition taking half done", slog.Any("srv", newSrv))
				return nil
			default:
				panic(step.ErrRootTypeUnexpected(curSem))
			}
		default:
			panic(state.ErrPolarityUnexpected(curSt))
		}
	default:
		panic(step.ErrTermTypeUnexpected(proc.Term))
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
	// Owner ID
	OID chnl.ID
	Ctx []chnl.ID
}

// Transition
type TranSpec struct {
	// Deal ID
	DID ID
	// Process ID
	PID chnl.ID
	// Agent Access Key
	Key  ak.ADT
	Term step.Term
}

// func (s *dealService) retrieveCtx() state.Context {
// }

// aka checkExp
func (s *dealService) checkProvider(
	ctx state.Context,
	zc state.ZC,
	t step.Term,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		if len(ctx.Linear) > 0 {
			return fmt.Errorf("context mismatch: want 0 items, got %v items", len(ctx.Linear))
		}
		// check via
		return state.CheckRoot(zc.C, state.OneRoot{})
	case step.WaitSpec:
		return step.ErrTermTypeMismatch(t, step.CloseSpec{})
	case step.SendSpec:
		// check via
		want, ok := zc.C.(state.TensorRoot)
		if !ok {
			return state.ErrRootTypeMismatch(zc.C, want)
		}
		// check value
		gotB, ok := ctx.Linear[term.B]
		if !ok {
			err := chnl.ErrMissingInCtx(term.B)
			s.log.Error("type checking failed",
				slog.Any("reason", err),
			)
			return err
		}
		err := state.CheckRoot(gotB, want.B)
		if err != nil {
			return err
		}
		// no cont to check
		return nil
	case step.RecvSpec:
		// check via
		want, ok := zc.C.(state.LolliRoot)
		if !ok {
			return state.ErrRootTypeMismatch(zc.C, want)
		}
		// check value
		gotY, ok := ctx.Linear[term.Y]
		if !ok {
			return chnl.ErrMissingInCtx(term.Y)
		}
		err := state.CheckRoot(gotY, want.Y)
		if err != nil {
			return err
		}
		// check cont
		return s.checkProvider(ctx, zc, term.Cont)
	case step.LabSpec:
		// check via
		want, ok := zc.C.(state.PlusRoot)
		if !ok {
			return state.ErrRootTypeMismatch(zc.C, want)
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
		want, ok := zc.C.(state.WithRoot)
		if !ok {
			return state.ErrRootTypeMismatch(zc.C, want)
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
			err := s.checkProvider(ctx, zc, gotCont)
			if err != nil {
				return err
			}
		}
		return nil
	case step.FwdSpec:
		if len(ctx.Linear) != 1 {
			return fmt.Errorf("context mismatch: want 1 item, got %v items", len(ctx.Linear))
		}
		gotD, ok := ctx.Linear[term.D]
		if !ok {
			return chnl.ErrMissingInCtx(term.D)
		}
		if zc.C.Pol() != gotD.Pol() {
			return state.ErrPolarityMismatch(zc.C, gotD)
		}
		return nil
	default:
		panic(step.ErrTermTypeUnexpected(t))
	}
}

func (s *dealService) checkClient(
	ctx state.Context,
	zc state.ZC,
	t step.Term,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		return step.ErrTermTypeMismatch(t, step.WaitSpec{})
	case step.WaitSpec:
		// check via
		gotX, ok := ctx.Linear[term.X]
		if !ok {
			return chnl.ErrMissingInCtx(term.X)
		}
		want, ok := gotX.(state.OneRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotX, want)
		}
		// check cont
		return s.checkClient(ctx, zc, term.Cont)
	case step.SendSpec:
		// check via
		gotA, ok := ctx.Linear[term.A]
		if !ok {
			return chnl.ErrMissingInCtx(term.A)
		}
		want, ok := gotA.(state.LolliRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotA, want)
		}
		// check value
		gotB, ok := ctx.Linear[term.B]
		if !ok {
			return chnl.ErrMissingInCtx(term.B)
		}
		err := state.CheckRoot(gotB, want.Y)
		if err != nil {
			return err
		}
		// no cont to check
		return nil
	case step.RecvSpec:
		// check via
		gotX, ok := ctx.Linear[term.X]
		if !ok {
			return chnl.ErrMissingInCtx(term.X)
		}
		want, ok := gotX.(state.TensorRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotX, want)
		}
		// check value
		gotY, ok := ctx.Linear[term.Y]
		if !ok {
			return chnl.ErrMissingInCtx(term.Y)
		}
		err := state.CheckRoot(gotY, want.B)
		if err != nil {
			return err
		}
		// check cont
		return s.checkClient(ctx, zc, term.Cont)
	case step.LabSpec:
		// check via
		gotC, ok := ctx.Linear[term.C]
		if !ok {
			return chnl.ErrMissingInCtx(term.C)
		}
		want, ok := gotC.(state.WithRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotC, want)
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
		gotZ, ok := ctx.Linear[term.Z]
		if !ok {
			return chnl.ErrMissingInCtx(term.Z)
		}
		want, ok := gotZ.(state.PlusRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotZ, want)
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
			err := s.checkClient(ctx, zc, gotCont)
			if err != nil {
				return err
			}
		}
		return nil
	case step.SpawnSpec:
		return nil
	case step.FwdSpec:
		if len(ctx.Linear) != 1 {
			return fmt.Errorf("context mismatch: want 1 item, got %v items", len(ctx.Linear))
		}
		gotC, ok := ctx.Linear[term.C]
		if !ok {
			return chnl.ErrMissingInCtx(term.C)
		}
		if gotC.Pol() != zc.C.Pol() {
			return state.ErrPolarityMismatch(gotC, zc.C)
		}
		return nil
	default:
		panic(step.ErrTermTypeUnexpected(t))
	}
}

func convertToCfg(chnls []chnl.Root) map[chnl.ID]chnl.Root {
	cfg := make(map[chnl.ID]chnl.Root, len(chnls))
	for _, ch := range chnls {
		cfg[ch.ID] = ch
	}
	return cfg
}

func join(chnls []chnl.Root, states map[state.ID]state.Root) state.Context {
	linear := make(map[core.Placeholder]state.Root, len(chnls))
	for _, ch := range chnls {
		linear[ch.ID] = states[ch.StID]
	}
	return state.Context{Linear: linear}
}
