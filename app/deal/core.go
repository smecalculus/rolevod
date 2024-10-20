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

type Environment struct {
	seats  map[seat.ID]seat.SeatRoot
	states map[state.ID]state.Root
}

func (e Environment) Contains(sid seat.ID) bool {
	_, ok := e.seats[sid]
	return ok
}

func (e Environment) LookupPE(sid seat.ID) state.PE {
	decl := e.seats[sid]
	return state.PE{Z: sym.New(decl.PE.Name), C: e.states[decl.PE.StID]}
}

func (e Environment) LookupCEs(sid seat.ID) []state.PE {
	decl := e.seats[sid]
	ctx := []state.PE{}
	for _, spec := range decl.CEs {
		ctx = append(ctx, state.PE{Z: sym.New(spec.Name), C: e.states[spec.StID]})
	}
	return ctx
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
	seats    seat.SeatRepo
	chnls    chnl.Repo
	steps    step.Repo
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newDealService(
	deals dealRepo,
	seats seat.SeatRepo,
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
	wantSpec, err := s.seats.SelectByID(gotSpec.Decl)
	if err != nil {
		s.log.Error("seat selection failed",
			slog.Any("reason", err),
			slog.Any("id", gotSpec.Decl),
		)
		return step.ProcRoot{}, err
	}
	newPE := chnl.Root{
		ID:   id.New(),
		Name: wantSpec.PE.Name,
		StID: wantSpec.PE.StID,
	}
	err = s.chnls.Insert(newPE)
	if err != nil {
		s.log.Error("providable endpoint insertion failed",
			slog.Any("reason", err),
			slog.Any("pe", newPE),
		)
		return step.ProcRoot{}, err
	}
	if len(gotSpec.Ctx) > 0 {
		err = s.chnls.Transfer(gotSpec.Owner, newPE.ID, gotSpec.Ctx)
		if err != nil {
			s.log.Error("context transfer failed",
				slog.Any("reason", err),
				slog.Any("from", gotSpec.Owner),
				slog.Any("to", newPE.ID),
				slog.Any("ctx", gotSpec.Ctx),
			)
			return step.ProcRoot{}, err
		}
	}
	newProc := step.ProcRoot{
		ID:  id.New(),
		PID: newPE.ID,
		Ctx: gotSpec.Ctx,
		Term: step.CTASpec{
			AK:  ak.New(),
			SID: gotSpec.Decl,
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
	curStep, err := s.steps.SelectByPID(spec.PID)
	if err != nil {
		s.log.Error("process selection failed",
			slog.Any("reason", err),
			slog.Any("pid", spec.PID),
		)
		return err
	}
	if curStep == nil {
		err = step.ErrDoesNotExist(spec.PID)
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
			slog.Any("pid", spec.PID),
		)
		return err
	}
	curProc, ok := curStep.(step.ProcRoot)
	if !ok {
		err = step.ErrRootTypeMismatch(curStep, step.ProcRoot{})
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
			slog.Any("pid", spec.PID),
		)
		return err
	}
	_, ok = curProc.Term.(step.CTASpec)
	if !ok {
		err = step.ErrTermTypeMismatch(spec.Term, step.CTASpec{})
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
			slog.Any("pid", curProc.PID),
		)
		return err
	}
	seatIDs := step.CollectEnv(spec.Term)
	seats, err := s.seats.SelectEnv(seatIDs)
	if err != nil {
		s.log.Error("seats selection failed",
			slog.Any("reason", err),
			slog.Any("pid", curProc.PID),
			slog.Any("ids", seatIDs),
		)
		return err
	}
	pe, err := s.chnls.SelectByID(curProc.PID)
	if err != nil {
		s.log.Error("providable endpoint selection failed",
			slog.Any("reason", err),
			slog.Any("pid", curProc.PID),
		)
		return err
	}
	ceIDs := step.CollectCtx(curProc.PID, spec.Term)
	ces, err := s.chnls.SelectCtx(curProc.PID, ceIDs)
	if err != nil {
		s.log.Error("consumable endpoints selection failed",
			slog.Any("reason", err),
			slog.Any("pid", curProc.PID),
			slog.Any("ids", ceIDs),
		)
		return err
	}
	envIDs := seat.CollectStIDs(maps.Values(seats))
	ctxIDs := chnl.CollectStIDs(append(ces, pe))
	states, err := s.states.SelectEnv(append(envIDs, ctxIDs...))
	if err != nil {
		s.log.Error("states selection failed",
			slog.Any("reason", err),
			slog.Any("env", envIDs),
			slog.Any("ctx", ctxIDs),
		)
		return err
	}
	env := Environment{seats, states}
	ctx := convertToCtx(ces, states)
	zc := state.PE{Z: pe.ID, C: states[pe.StID]}
	// type checking
	err = s.checkState(env, ctx, zc, spec.Term)
	if err != nil {
		s.log.Error("transition taking failed", slog.Any("reason", err))
		return err
	}
	// step taking
	cfg := convertToCfg(append(ces, pe))
	curProc.Term = spec.Term
	return s.takeProcWith(curProc, cfg, states)
}

func (s *dealService) takeProc(
	proc step.ProcRoot,
) (err error) {
	s.log.Debug("transition taking started", slog.Any("proc", proc))
	pe, err := s.chnls.SelectByID(proc.PID)
	if err != nil {
		s.log.Error("providable endpoint selection failed",
			slog.Any("reason", err),
			slog.Any("pid", proc.PID),
		)
		return err
	}
	ces, err := s.chnls.SelectCtx(proc.PID, proc.Ctx)
	if err != nil {
		s.log.Error("consumable endpoints selection failed",
			slog.Any("reason", err),
			slog.Any("ctx", proc.Ctx),
		)
		return err
	}
	stIDs := chnl.CollectStIDs(append(ces, pe))
	env, err := s.states.SelectEnv(stIDs)
	if err != nil {
		s.log.Error("env selection failed",
			slog.Any("reason", err),
			slog.Any("ids", stIDs),
		)
		return err
	}
	cfg := convertToCfg(append(ces, pe))
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
		newProc, err := s.Involve(PartSpec{Decl: term.SeatID, Owner: proc.PID, Ctx: term.Ctx})
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
	Deal  ID
	Decl  seat.ID
	Owner chnl.ID
	Ctx   []chnl.ID
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

func (s *dealService) checkState(
	env Environment,
	ctx state.Context,
	pe state.PE,
	t step.Term,
) error {
	if pe.Z == t.Via() {
		return s.checkProvider(env, ctx, pe, t)
	}
	return s.checkClient(env, ctx, pe, t)
}

// aka checkExp
func (s *dealService) checkProvider(
	env Environment,
	ctx state.Context,
	pe state.PE,
	t step.Term,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		if len(ctx.Linear) > 0 {
			err := fmt.Errorf("context mismatch: want 0 items, got %v items", len(ctx.Linear))
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check via
		return state.CheckRoot(pe.C, state.OneRoot{})
	case step.WaitSpec:
		err := step.ErrTermTypeMismatch(t, step.CloseSpec{})
		s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
		return err
	case step.SendSpec:
		// check via
		want, ok := pe.C.(state.TensorRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(pe.C, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check value
		gotB, ok := ctx.Linear[term.B]
		if !ok {
			err := chnl.ErrMissingInCtx(term.B)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		err := state.CheckRoot(gotB, want.B)
		if err != nil {
			return err
		}
		// no cont to check
		delete(ctx.Linear, term.B)
		pe.C = want.C
		return nil
	case step.RecvSpec:
		// check via
		want, ok := pe.C.(state.LolliRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(pe.C, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check value
		gotY, ok := ctx.Linear[term.Y]
		if !ok {
			err := chnl.ErrMissingInCtx(term.Y)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		err := state.CheckRoot(gotY, want.Y)
		if err != nil {
			return err
		}
		// check cont
		ctx.Linear[term.Y] = want.Y
		pe.C = want.Z
		return s.checkState(env, ctx, pe, term.Cont)
	case step.LabSpec:
		// check via
		want, ok := pe.C.(state.PlusRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(pe.C, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check label
		c, ok := want.Choices[term.L]
		if !ok {
			err := fmt.Errorf("label mismatch: want %q, got nothing", term.L)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// no cont to check
		pe.C = c
		return nil
	case step.CaseSpec:
		// check via
		want, ok := pe.C.(state.WithRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(pe.C, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check conts
		if len(term.Conts) != len(want.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(want.Choices), len(term.Conts))
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		for wantL, wantC := range want.Choices {
			gotC, ok := term.Conts[wantL]
			if !ok {
				err := fmt.Errorf("label mismatch: want %q, got nothing", wantL)
				s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
				return err
			}
			pe.C = wantC
			err := s.checkState(env, ctx, pe, gotC)
			if err != nil {
				return err
			}
		}
		return nil
	case step.FwdSpec:
		if len(ctx.Linear) != 1 {
			err := fmt.Errorf("context mismatch: want 1 item, got %v items", len(ctx.Linear))
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		gotD, ok := ctx.Linear[term.D]
		if !ok {
			err := chnl.ErrMissingInCtx(term.D)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		if gotD.Pol() != pe.C.Pol() {
			err := state.ErrPolarityMismatch(gotD, pe.C)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		return state.CheckRoot(gotD, pe.C)
	default:
		panic(step.ErrTermTypeUnexpected(t))
	}
}

func (s *dealService) checkClient(
	env Environment,
	ctx state.Context,
	pe state.PE,
	t step.Term,
) error {
	switch got := t.(type) {
	case step.CloseSpec:
		err := step.ErrTermTypeMismatch(t, step.WaitSpec{})
		s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
		return err
	case step.WaitSpec:
		// check via
		gotX, ok := ctx.Linear[got.X]
		if !ok {
			err := chnl.ErrMissingInCtx(got.X)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		want, ok := gotX.(state.OneRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotX, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check cont
		delete(ctx.Linear, got.X)
		return s.checkState(env, ctx, pe, got.Cont)
	case step.SendSpec:
		// check via
		gotA, ok := ctx.Linear[got.A]
		if !ok {
			err := chnl.ErrMissingInCtx(got.A)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		want, ok := gotA.(state.LolliRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotA, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check value
		gotB, ok := ctx.Linear[got.B]
		if !ok {
			err := chnl.ErrMissingInCtx(got.B)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		err := state.CheckRoot(gotB, want.Y)
		if err != nil {
			return err
		}
		// no cont to check
		delete(ctx.Linear, got.B)
		ctx.Linear[got.A] = want.Z
		return nil
	case step.RecvSpec:
		// check via
		gotX, ok := ctx.Linear[got.X]
		if !ok {
			err := chnl.ErrMissingInCtx(got.X)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		want, ok := gotX.(state.TensorRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotX, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check value
		gotY, ok := ctx.Linear[got.Y]
		if !ok {
			err := chnl.ErrMissingInCtx(got.Y)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		err := state.CheckRoot(gotY, want.B)
		if err != nil {
			return err
		}
		// check cont
		ctx.Linear[got.Y] = want.B
		pe.C = want.C
		return s.checkState(env, ctx, pe, got.Cont)
	case step.LabSpec:
		// check via
		gotC, ok := ctx.Linear[got.C]
		if !ok {
			err := chnl.ErrMissingInCtx(got.C)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		want, ok := gotC.(state.WithRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotC, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check label
		c, ok := want.Choices[got.L]
		if !ok {
			err := fmt.Errorf("label mismatch: want %q, got nothing", got.L)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// no cont to check
		ctx.Linear[got.C] = c
		return nil
	case step.CaseSpec:
		// check via
		gotZ, ok := ctx.Linear[got.Z]
		if !ok {
			err := chnl.ErrMissingInCtx(got.Z)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		want, ok := gotZ.(state.PlusRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotZ, want)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check conts
		if len(got.Conts) != len(want.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(want.Choices), len(got.Conts))
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		for wantL, wantC := range want.Choices {
			gotC, ok := got.Conts[wantL]
			if !ok {
				err := fmt.Errorf("label mismatch: want %q, got nothing", wantL)
				s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
				return err
			}
			ctx.Linear[got.Z] = wantC
			err := s.checkState(env, ctx, pe, gotC)
			if err != nil {
				return err
			}
		}
		return nil
	case step.SpawnSpec:
		if !env.Contains(got.SeatID) {
			err := seat.ErrRootMissingInEnv(got.SeatID)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		wantCtx := env.LookupCEs(got.SeatID)
		if len(got.Ctx) != len(wantCtx) {
			err := fmt.Errorf("context mismatch: want %v items, got %v items", len(wantCtx), len(got.Ctx))
			s.log.Error("type checking failed",
				slog.Any("reason", err),
				slog.Any("via", t.Via()),
				slog.Any("want", wantCtx),
				slog.Any("got", got.Ctx),
			)
			return err
		}
		if len(got.Ctx) == 0 {
			return nil
		}
		for i, gotID := range got.Ctx {
			gotSt := ctx.Linear[gotID]
			err := state.CheckRoot(gotSt, wantCtx[i].C)
			if err != nil {
				s.log.Error("type checking failed",
					slog.Any("reason", err),
					slog.Any("via", t.Via()),
					slog.Any("want", wantCtx[i]),
					slog.Any("got", gotID),
				)
				return err
			}
			delete(ctx.Linear, gotID)
		}
		ctx.Linear[got.Z] = env.LookupPE(got.SeatID).C
		return s.checkState(env, ctx, pe, got.Cont)
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

func convertToCtx(chnls []chnl.Root, states map[state.ID]state.Root) state.Context {
	linear := make(map[core.Placeholder]state.Root, len(chnls))
	for _, ch := range chnls {
		linear[ch.ID] = states[ch.StID]
	}
	return state.Context{Linear: linear}
}
