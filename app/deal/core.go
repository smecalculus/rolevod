package deal

import (
	"fmt"
	"log/slog"

	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

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

func (s *dealService) Involve(spec PartSpec) (step.ProcRoot, error) {
	s.log.Debug("seat involvement started", slog.Any("spec", spec))
	procSpec, err := s.seats.Retrieve(spec.SeatID)
	if err != nil {
		s.log.Error("seat selection failed",
			slog.Any("reason", err),
			slog.Any("id", spec.SeatID),
		)
		return step.ProcRoot{}, err
	}
	if len(spec.Ctx) != len(procSpec.Ctx) {
		err := fmt.Errorf("context mismatch: want %v items, got %v items", len(procSpec.Ctx), len(spec.Ctx))
		s.log.Error("seat involvement failed",
			slog.Any("reason", err),
			slog.Any("ctx", spec.Ctx),
		)
		return step.ProcRoot{}, err
	}
	procVia := chnl.Root{
		ID:   id.New(),
		Name: procSpec.Via.Name,
		StID: procSpec.Via.StID,
		St:   procSpec.Via.St,
	}
	err = s.chnls.Insert(procVia)
	if err != nil {
		s.log.Error("via insertion failed",
			slog.Any("reason", err),
			slog.Any("via", procVia),
		)
		return step.ProcRoot{}, err
	}
	if len(spec.Ctx) > 0 {
		curCtx, err := s.chnls.SelectCtx(spec.OID, maps.Values(spec.Ctx))
		if err != nil {
			s.log.Error("context selection failed",
				slog.Any("reason", err),
				slog.Any("ctx", spec.Ctx),
			)
			return step.ProcRoot{}, err
		}
		for i, got := range curCtx {
			// TODO обеспечить порядок
			// TODO проверять по значению, а не по ссылке
			err = checkState(got.St, procSpec.Ctx[i].St)
			if err != nil {
				s.log.Error("type checking failed",
					slog.Any("reason", err),
					slog.Any("st", got.St),
				)
				return step.ProcRoot{}, err
			}
		}
		err = s.chnls.Transfer(spec.OID, procVia.ID, maps.Values(spec.Ctx))
		if err != nil {
			s.log.Error("context transfer failed",
				slog.Any("reason", err),
				slog.Any("from", spec.OID),
				slog.Any("to", procVia.ID),
				slog.Any("ctx", spec.Ctx),
			)
			return step.ProcRoot{}, err
		}
	}
	newProc := step.ProcRoot{
		ID:  id.New(),
		PID: procVia.ID,
		Ctx: spec.Ctx,
		Term: step.CTASpec{
			AK:  ak.New(),
			SID: spec.SeatID,
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
	zc, err := s.chnls.SelectByID(proc.PID)
	if err != nil {
		s.log.Error("zc selection failed",
			slog.Any("reason", err),
			slog.Any("id", proc.PID),
		)
		return err
	}
	chnls, err := s.chnls.SelectCtx(proc.PID, maps.Values(proc.Ctx))
	if err != nil {
		s.log.Error("context selection failed",
			slog.Any("reason", err),
			slog.Any("ctx", proc.Ctx),
		)
		return err
	}
	// type checking
	stIDs := chnl.CollectStIDs(append(chnls, zc))
	env, err := s.states.SelectEnv(stIDs)
	if err != nil {
		s.log.Error("env selection failed",
			slog.Any("reason", err),
			slog.Any("ids", stIDs),
		)
		return err
	}
	ctx := convertToCtx(chnls)
	cfg := convertToCfg(append(chnls, zc))
	if spec.Term.Via() == spec.PID {
		s.checkProvider(env, ctx, cfg, spec.Term)
	} else {
		s.checkClient(env, ctx, cfg, spec.Term)
	}
	// step taking
	proc.Term = spec.Term
	return s.takeProcWith(proc, cfg, env)
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
	ctx, err := s.chnls.SelectCtx(proc.PID, maps.Values(proc.Ctx))
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
			newProc = step.ProcRoot{
				ID:   id.New(),
				PID:  proc.PID,
				Term: term.Cont,
			}
		case step.FwdSpec:
			c, ok := cfg[val.C.(chnl.ID)]
			if !ok {
				err = chnl.ErrMissingInCfg(c.ID)
				s.log.Error("transition taking failed",
					slog.Any("reason", err),
				)
				return err
			}
			err := s.chnls.Transfer(msg.PID, msg.PID, []chnl.ID{c.ID})
			if err != nil {
				s.log.Error("channel transfer failed",
					slog.Any("reason", err),
					slog.Any("id", c.ID),
				)
				return err
			}
			newProc = step.ProcRoot{
				ID:   id.New(),
				PID:  proc.PID,
				Ctx:  map[string]id.ADT{c.Name: c.ID},
				Term: step.Subst(term, val.D, c.ID),
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
			Ctx:  map[string]id.ADT{b.Name: b.ID},
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
			err = state.ErrDoesNotExist(viaCh.StID)
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
					Ctx:  map[string]chnl.ID{c.Name: c.ID},
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
					Ctx:  map[string]chnl.ID{d.Name: d.ID},
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
					Ctx:  map[string]chnl.ID{d.Name: d.ID},
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
					Ctx:  map[string]chnl.ID{c.Name: c.ID},
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
			panic(state.ErrUnexpectedPolarity(curSt.Pol()))
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
	Ctx map[chnl.Name]chnl.ID
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

// aka checkExp
func (s *dealService) checkProvider(
	env map[state.ID]state.Root,
	ctx map[chnl.Name]chnl.Root,
	cfg map[chnl.ID]chnl.Root,
	t step.Term,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		// check via
		gotA, err := findState(env, ctx, cfg, term.A)
		if err != nil {
			return err
		}
		return checkProvider(gotA, state.OneRoot{})
	case step.WaitSpec:
		// check via
		gotX, err := findState(env, ctx, cfg, term.X)
		if err != nil {
			return err
		}
		_, ok := gotX.(state.OneRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotX, state.OneRoot{})
		}
		// check cont
		return s.checkProvider(env, ctx, cfg, term.Cont)
	case step.SendSpec:
		// check via
		gotA, err := findState(env, ctx, cfg, term.A)
		if err != nil {
			return err
		}
		want, ok := gotA.(state.TensorRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotA, state.TensorRoot{})
		}
		// check value
		gotB, err := findState(env, ctx, cfg, term.B)
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
		gotX, err := findState(env, ctx, cfg, term.X)
		if err != nil {
			return err
		}
		want, ok := gotX.(state.LolliRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotX, state.LolliRoot{})
		}
		// check value
		gotY, err := findState(env, ctx, cfg, term.Y)
		if err != nil {
			return err
		}
		err = checkProvider(gotY, want.Y)
		if err != nil {
			return err
		}
		// check cont
		return s.checkProvider(env, ctx, cfg, term.Cont)
	case step.LabSpec:
		// check via
		gotC, err := findState(env, ctx, cfg, term.C)
		if err != nil {
			return err
		}
		want, ok := gotC.(state.PlusRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotC, state.PlusRoot{})
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
		gotZ, err := findState(env, ctx, cfg, term.Z)
		if err != nil {
			return err
		}
		want, ok := gotZ.(state.WithRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotZ, state.WithRoot{})
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
			err := s.checkProvider(env, ctx, cfg, gotCont)
			if err != nil {
				return err
			}
		}
		return nil
	case step.SpawnSpec:
		return nil
	case step.FwdSpec:
		gotC, err := findState(env, ctx, cfg, term.C)
		if err != nil {
			return err
		}
		gotD, err := findState(env, ctx, cfg, term.D)
		if err != nil {
			return err
		}
		if gotC.Pol() != gotD.Pol() {
			return fmt.Errorf("polarity mismatch: want C==D, got %v!=%v", gotC.Pol(), gotD.Pol())
		}
		return nil
	default:
		panic(step.ErrTermTypeUnexpected(t))
	}
}

func (s *dealService) checkClient(
	env map[state.ID]state.Root,
	ctx map[chnl.Name]chnl.Root,
	cfg map[chnl.ID]chnl.Root,
	t step.Term,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		// check via
		gotA, err := findState(env, ctx, cfg, term.A)
		if err != nil {
			return err
		}
		return checkClient(gotA, state.OneRoot{})
	case step.WaitSpec:
		// check via
		gotX, err := findState(env, ctx, cfg, term.X)
		if err != nil {
			return err
		}
		_, ok := gotX.(state.OneRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotX, state.OneRoot{})
		}
		// check cont
		return s.checkClient(env, ctx, cfg, term.Cont)
	case step.SendSpec:
		// check via
		gotA, err := findState(env, ctx, cfg, term.A)
		if err != nil {
			return err
		}
		want, ok := gotA.(state.LolliRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotA, state.LolliRoot{})
		}
		// check value
		gotB, err := findState(env, ctx, cfg, term.B)
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
		gotX, err := findState(env, ctx, cfg, term.X)
		if err != nil {
			return err
		}
		want, ok := gotX.(state.TensorRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotX, state.TensorRoot{})
		}
		// check value
		gotY, err := findState(env, ctx, cfg, term.Y)
		if err != nil {
			return err
		}
		err = checkProvider(gotY, want.B)
		if err != nil {
			return err
		}
		// check cont
		return s.checkClient(env, ctx, cfg, term.Cont)
	case step.LabSpec:
		// check via
		gotC, err := findState(env, ctx, cfg, term.C)
		if err != nil {
			return err
		}
		want, ok := gotC.(state.WithRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotC, state.WithRoot{})
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
		gotZ, err := findState(env, ctx, cfg, term.Z)
		if err != nil {
			return err
		}
		want, ok := gotZ.(state.PlusRoot)
		if !ok {
			return state.ErrRootTypeMismatch(gotZ, state.PlusRoot{})
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
			err := s.checkClient(env, ctx, cfg, gotCont)
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
		panic(step.ErrTermTypeUnexpected(t))
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
			return state.ErrRootTypeMismatch(got, want)
		}
		return nil
	case state.TensorRoot:
		gotSt, ok := got.(state.TensorRoot)
		if !ok {
			return state.ErrRootTypeMismatch(got, want)
		}
		err := checkProvider(gotSt.B, wantSt.B)
		if err != nil {
			return err
		}
		return checkProvider(gotSt.C, wantSt.C)
	case state.LolliRoot:
		gotSt, ok := got.(state.LolliRoot)
		if !ok {
			return state.ErrRootTypeMismatch(got, want)
		}
		err := checkProvider(gotSt.Y, wantSt.Y)
		if err != nil {
			return err
		}
		return checkProvider(gotSt.Z, wantSt.Z)
	case state.PlusRoot:
		gotSt, ok := got.(state.PlusRoot)
		if !ok {
			return state.ErrRootTypeMismatch(got, want)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantSt.Choices), len(gotSt.Choices))
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
			return state.ErrRootTypeMismatch(got, want)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantSt.Choices), len(gotSt.Choices))
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
		panic(state.ErrRootTypeUnexpected(want))
	}
}

func checkClient(got, want state.Root) error {
	switch wantSt := want.(type) {
	case state.OneRoot:
		_, ok := got.(state.OneRoot)
		if !ok {
			return state.ErrRootTypeMismatch(got, want)
		}
		return nil
	case state.TensorRoot:
		gotSt, ok := got.(state.LolliRoot)
		if !ok {
			return state.ErrRootTypeMismatch(got, want)
		}
		err := checkClient(gotSt.Y, wantSt.B)
		if err != nil {
			return err
		}
		return checkClient(gotSt.Z, wantSt.C)
	case state.LolliRoot:
		gotSt, ok := got.(state.TensorRoot)
		if !ok {
			return state.ErrRootTypeMismatch(got, want)
		}
		err := checkClient(gotSt.B, wantSt.Y)
		if err != nil {
			return err
		}
		return checkClient(gotSt.C, wantSt.Z)
	case state.PlusRoot:
		gotSt, ok := got.(state.WithRoot)
		if !ok {
			return state.ErrRootTypeMismatch(got, want)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantSt.Choices), len(gotSt.Choices))
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
			return state.ErrRootTypeMismatch(got, want)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantSt.Choices), len(gotSt.Choices))
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
		panic(state.ErrRootTypeUnexpected(want))
	}
}

func findState(
	env map[state.ID]state.Root,
	ctx map[chnl.Name]chnl.Root,
	cfg map[chnl.ID]chnl.Root,
	ph core.Placeholder,
) (state.Root, error) {
	var gotCh chnl.Root
	switch val := ph.(type) {
	case chnl.ID:
		got, ok := cfg[val]
		if !ok {
			return nil, chnl.ErrMissingInCfg(val)
		}
		gotCh = got
	case sym.ADT:
		got, ok := ctx[val.Name()]
		if !ok {
			return nil, chnl.ErrMissingInCtx(val)
		}
		gotCh = got
	}
	if gotCh.St == nil {
		return nil, chnl.ErrAlreadyClosed(gotCh.ID)
	}
	gotSt, ok := env[gotCh.StID]
	if !ok {
		return nil, state.ErrDoesNotExist(gotCh.StID)
	}
	return gotSt, nil
}

func convertToCfg(chnls []chnl.Root) map[chnl.ID]chnl.Root {
	cfg := make(map[chnl.ID]chnl.Root, len(chnls))
	for _, ch := range chnls {
		cfg[ch.ID] = ch
	}
	return cfg
}

func convertToCtx(chnls []chnl.Root) map[chnl.Name]chnl.Root {
	ctx := make(map[chnl.Name]chnl.Root, len(chnls))
	for _, ch := range chnls {
		ctx[ch.Name] = ch
	}
	return ctx
}
