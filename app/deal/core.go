package deal

import (
	"fmt"
	"log/slog"

	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/pol"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/sig"
)

type Name = string
type ID = id.ADT

type Spec struct {
	Name Name
}

type Ref struct {
	ID   ID
	Name Name
}

// aka Configuration or Eta
type Root struct {
	ID       ID
	Name     Name
	Children []Ref
	Sigs     []sig.Ref
}

type Environment struct {
	sigs   map[sig.ID]sig.Root
	roles  map[role.ID]role.Root
	states map[state.ID]state.Root
}

func (e Environment) Contains(id sig.ID) bool {
	_, ok := e.sigs[id]
	return ok
}

func (e Environment) LookupPE(id sig.ID) state.EP {
	decl := e.sigs[id]
	role := e.roles[decl.PE.Role]
	return state.EP{Z: sym.New(decl.PE.Name), C: e.states[role.StateID]}
}

func (e Environment) LookupCEs(id sig.ID) []state.EP {
	decl := e.sigs[id]
	ces := []state.EP{}
	for _, ce := range decl.CEs {
		role := e.roles[decl.PE.Role]
		ces = append(ces, state.EP{Z: sym.New(ce.Name), C: e.states[role.StateID]})
	}
	return ces
}

type Configuration struct {
	chnls  map[chnl.ID]chnl.Root
	states map[state.ID]state.Root
}

func (c *Configuration) LookupCh(id chnl.ID) (chnl.Root, bool) {
	ch, ok := c.chnls[id]
	if !ok {
		return chnl.Root{}, false
	}
	return ch, true
}

func (c *Configuration) LookupSt(id chnl.ID) (state.Root, bool) {
	ch, ok := c.chnls[id]
	if !ok {
		return nil, false
	}
	if ch.StateID == nil {
		return nil, false
	}
	st, ok := c.states[*ch.StateID]
	if !ok {
		panic(state.ErrMissingInCfg(*ch.StateID))
	}
	return st, true
}

func (c *Configuration) Add(ch chnl.Root) {
	c.chnls[ch.ID] = ch
}

func (c *Configuration) Remove(ids ...chnl.ID) {
	for _, id := range ids {
		delete(c.chnls, id)
	}
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)

type API interface {
	Create(Spec) (Root, error)
	Retrieve(ID) (Root, error)
	RetreiveAll() ([]Ref, error)
	Establish(KinshipSpec) error
	Involve(PartSpec) (chnl.Root, error)
	Take(TranSpec) error
}

type service struct {
	deals    repo
	roles    role.Repo
	sigs     sig.Repo
	chnls    chnl.Repo
	steps    step.Repo
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func newService(
	deals repo,
	roles role.Repo,
	sigs sig.Repo,
	chnls chnl.Repo,
	steps step.Repo,
	states state.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *service {
	name := slog.String("name", "dealService")
	return &service{
		deals, roles, sigs, chnls, steps, states, kinships, l.With(name),
	}
}

func (s *service) Create(spec Spec) (Root, error) {
	s.log.Debug("deal creation started", slog.Any("spec", spec))
	root := Root{
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

func (s *service) Retrieve(id ID) (Root, error) {
	root, err := s.deals.SelectByID(id)
	if err != nil {
		return Root{}, err
	}
	root.Children, err = s.deals.SelectChildren(id)
	if err != nil {
		return Root{}, err
	}
	return root, nil
}

func (s *service) RetreiveAll() ([]Ref, error) {
	return s.deals.SelectAll()
}

func (s *service) Establish(spec KinshipSpec) error {
	s.log.Debug("kinship establishment started", slog.Any("spec", spec))
	var children []Ref
	for _, id := range spec.ChildIDs {
		children = append(children, Ref{ID: id})
	}
	root := KinshipRoot{
		Parent:   Ref{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinships.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("kinship establishment succeeded", slog.Any("root", root))
	return nil
}

func (s *service) Involve(gotSpec PartSpec) (chnl.Root, error) {
	s.log.Debug("sig involvement started", slog.Any("spec", gotSpec))
	wantSig, err := s.sigs.SelectByID(gotSpec.Sig)
	if err != nil {
		s.log.Error("signature selection failed",
			slog.Any("reason", err),
			slog.Any("id", gotSpec.Sig),
		)
		return chnl.Root{}, err
	}
	wantRole, err := s.roles.SelectByID(wantSig.PE.Role)
	if err != nil {
		s.log.Error("role selection failed",
			slog.Any("reason", err),
			slog.Any("id", wantSig.PE.Role),
		)
		return chnl.Root{}, err
	}
	newPE := chnl.Root{
		ID:      id.New(),
		Name:    wantSig.PE.Name,
		StateID: &wantRole.StateID,
	}
	err = s.chnls.Insert(newPE)
	if err != nil {
		s.log.Error("providable endpoint insertion failed",
			slog.Any("reason", err),
			slog.Any("pe", newPE),
		)
		return chnl.Root{}, err
	}
	if len(gotSpec.TEs) > 0 {
		err = s.chnls.Transfer(gotSpec.Owner, newPE.ID, gotSpec.TEs)
		if err != nil {
			s.log.Error("context transfer failed",
				slog.Any("reason", err),
				slog.Any("from", gotSpec.Owner),
				slog.Any("to", newPE.ID),
				slog.Any("tes", gotSpec.TEs),
			)
			return chnl.Root{}, err
		}
	}
	newProc := step.ProcRoot{
		ID:  id.New(),
		PID: newPE.ID,
		Term: step.CTASpec{
			AK:  ak.New(),
			Sig: gotSpec.Sig,
		},
	}
	err = s.steps.Insert(newProc)
	if err != nil {
		s.log.Error("process insertion failed",
			slog.Any("reason", err),
			slog.Any("proc", newProc),
		)
		return chnl.Root{}, err
	}
	s.log.Debug("sig involvement succeeded", slog.Any("proc", newProc))
	return newPE, nil
}

func (s *service) Take(spec TranSpec) error {
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
	proc, ok := curStep.(step.ProcRoot)
	if !ok {
		err = step.ErrRootTypeMismatch(curStep, step.ProcRoot{})
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
			slog.Any("pid", spec.PID),
		)
		return err
	}
	_, ok = proc.Term.(step.CTASpec)
	if !ok {
		err = step.ErrTermTypeMismatch(spec.Term, step.CTASpec{})
		s.log.Error("transition taking failed",
			slog.Any("reason", err),
			slog.Any("pid", proc.PID),
		)
		return err
	}
	sigIDs := step.CollectEnv(spec.Term)
	sigs, err := s.sigs.SelectEnv(sigIDs)
	if err != nil {
		s.log.Error("signatures selection failed",
			slog.Any("reason", err),
			slog.Any("pid", proc.PID),
			slog.Any("ids", sigIDs),
		)
		return err
	}
	roleIDs := sig.CollectEnv(maps.Values(sigs))
	roles, err := s.roles.SelectEnv(roleIDs)
	if err != nil {
		s.log.Error("roles selection failed",
			slog.Any("reason", err),
			slog.Any("pid", proc.PID),
			slog.Any("ids", roleIDs),
		)
		return err
	}
	pe, err := s.chnls.SelectByID(proc.PID)
	if err != nil {
		s.log.Error("providable endpoint selection failed",
			slog.Any("reason", err),
			slog.Any("id", proc.PID),
		)
		return err
	}
	ceIDs := step.CollectCtx(proc.PID, spec.Term)
	ces, err := s.chnls.SelectCtx(proc.PID, ceIDs)
	if err != nil {
		s.log.Error("consumable endpoints selection failed",
			slog.Any("reason", err),
			slog.Any("pid", proc.PID),
			slog.Any("ids", ceIDs),
		)
		return err
	}
	envIDs := role.CollectEnv(maps.Values(roles))
	ctxIDs := chnl.CollectCtx(append(ces, pe))
	states, err := s.states.SelectEnv(append(envIDs, ctxIDs...))
	if err != nil {
		s.log.Error("states selection failed",
			slog.Any("reason", err),
			slog.Any("env", envIDs),
			slog.Any("ctx", ctxIDs),
		)
		return err
	}
	env := Environment{sigs, roles, states}
	ctx := convertToCtx(ces, states)
	zc := state.EP{Z: pe.ID, C: states[*pe.StateID]}
	// type checking
	err = s.checkState(env, ctx, zc, spec.Term)
	if err != nil {
		s.log.Error("transition taking failed", slog.Any("reason", err))
		return err
	}
	// step taking
	cfg := Configuration{chnls: convertToCfg(append(ces, pe)), states: states}
	proc.Term = spec.Term
	return s.takeProcWith(proc, cfg)
}

func (s *service) takeProc(
	proc step.ProcRoot,
) (err error) {
	s.log.Debug("transition taking started", slog.Any("proc", proc))
	pe, err := s.chnls.SelectByID(proc.PID)
	if err != nil {
		s.log.Error("providable endpoint selection failed",
			slog.Any("reason", err),
			slog.Any("id", proc.PID),
		)
		return err
	}
	ceIDs := step.CollectCtx(proc.PID, proc.Term)
	ces, err := s.chnls.SelectCtx(proc.PID, ceIDs)
	if err != nil {
		s.log.Error("consumable endpoints selection failed",
			slog.Any("reason", err),
			slog.Any("ids", ceIDs),
		)
		return err
	}
	stIDs := chnl.CollectCtx(append(ces, pe))
	states, err := s.states.SelectEnv(stIDs)
	if err != nil {
		s.log.Error("states selection failed",
			slog.Any("reason", err),
			slog.Any("ids", stIDs),
		)
		return err
	}
	cfg := Configuration{chnls: convertToCfg(append(ces, pe)), states: states}
	return s.takeProcWith(proc, cfg)
}

func (s *service) takeProcWith(
	proc step.ProcRoot,
	cfg Configuration,
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
		curVia, ok := cfg.LookupCh(viaID)
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
				slog.Any("vid", curVia.ID),
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
			err = fmt.Errorf("cont type mismatch: want %T, got %T", wait, srv.Cont)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("cont", srv.Cont),
			)
			return err
		}
		// consume and close channel
		finVia := chnl.Root{
			ID:      id.New(),
			Name:    curVia.Name,
			PreID:   &curVia.ID,
			StateID: nil,
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
		curVia, ok := cfg.LookupCh(viaID)
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
				slog.Any("vid", curVia.ID),
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
				// TODO log some ID
			)
			return err
		}
		var newProc step.ProcRoot
		switch val := msg.Val.(type) {
		case step.CloseSpec:
			// consume and close channel
			finVia := chnl.Root{
				ID:      id.New(),
				Name:    curVia.Name,
				PreID:   &curVia.ID,
				StateID: nil,
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
					slog.Any("from", msg.PID),
					slog.Any("to", proc.PID),
					slog.Any("id", d),
				)
				return err
			}
			newProc = step.ProcRoot{
				ID:   id.New(),
				PID:  proc.PID,
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
		curVia, ok := cfg.LookupCh(viaID)
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
				slog.Any("vid", curVia.ID),
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
				// TODO log some ID
			)
			return err
		}
		recv, ok := srv.Cont.(step.RecvSpec)
		if !ok {
			err = fmt.Errorf("cont type mismatch: want %T, got %T", recv, srv.Cont)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("cont", srv.Cont),
			)
			return err
		}
		curSt, ok := cfg.LookupSt(curVia.ID)
		if !ok {
			err = chnl.ErrMissingInCfg(curVia.ID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		nextID := curSt.(state.Prod).Next()
		newVia := chnl.Root{
			ID:      id.New(),
			Name:    curVia.Name,
			PreID:   &curVia.ID,
			StateID: &nextID,
		}
		err = s.chnls.Insert(newVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", newVia),
			)
			return err
		}
		bID, ok := term.B.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.B)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		b, ok := cfg.LookupCh(bID)
		if !ok {
			err = chnl.ErrMissingInCfg(bID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		err = s.chnls.Transfer(proc.PID, srv.PID, []chnl.ID{b.ID})
		if err != nil {
			s.log.Error("channel transfer failed",
				slog.Any("reason", err),
				slog.Any("from", proc.PID),
				slog.Any("to", srv.PID),
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
		curVia, ok := cfg.LookupCh(viaID)
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
				slog.Any("vid", curVia.ID),
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
				// TODO log some ID
			)
			return err
		}
		send, ok := msg.Val.(step.SendSpec)
		if !ok {
			err = fmt.Errorf("val type mismatch: want %T, got %T", send, msg.Val)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("val", msg.Val),
			)
			return err
		}
		curSt, ok := cfg.LookupSt(curVia.ID)
		if !ok {
			err = chnl.ErrMissingInCfg(curVia.ID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		nextID := curSt.(state.Prod).Next()
		newVia := chnl.Root{
			ID:      id.New(),
			Name:    curVia.Name,
			PreID:   &curVia.ID,
			StateID: &nextID,
		}
		err = s.chnls.Insert(newVia)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("via", newVia),
			)
			return err
		}
		bID, ok := send.B.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(send.B)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		b, ok := cfg.LookupCh(bID)
		if !ok {
			err = chnl.ErrMissingInCfg(bID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		err = s.chnls.Transfer(msg.PID, proc.PID, []chnl.ID{b.ID})
		if err != nil {
			s.log.Error("channel transfer failed",
				slog.Any("reason", err),
				slog.Any("from", msg.PID),
				slog.Any("to", proc.PID),
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
		viaID, ok := term.A.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.A)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg.LookupCh(viaID)
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
				slog.Any("vid", curVia.ID),
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
				// TODO log some ID
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
		curSt, ok := cfg.LookupSt(curVia.ID)
		if !ok {
			err = chnl.ErrMissingInCfg(curVia.ID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		nextID := curSt.(state.Sum).Next(term.L)
		newVia := chnl.Root{
			ID:      id.New(),
			Name:    curVia.Name,
			PreID:   &curVia.ID,
			StateID: &nextID,
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
			Term: step.Subst(cont.Conts[term.L], cont.X, newVia.ID),
		}
		return s.takeProc(newProc)
	case step.CaseSpec:
		viaID, ok := term.X.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.X)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg.LookupCh(viaID)
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
				slog.Any("vid", curVia.ID),
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
				// TODO log some ID
			)
			return err
		}
		lab, ok := msg.Val.(step.LabSpec)
		if !ok {
			err = fmt.Errorf("val type mismatch: want %T, got %T", lab, msg.Val)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
				slog.Any("val", msg.Val),
			)
			return err
		}
		curSt, ok := cfg.LookupSt(curVia.ID)
		if !ok {
			err = chnl.ErrMissingInCfg(curVia.ID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		nextID := curSt.(state.Sum).Next(lab.L)
		newVia := chnl.Root{
			ID:      id.New(),
			Name:    curVia.Name,
			PreID:   &curVia.ID,
			StateID: &nextID,
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
			Term: step.Subst(term.Conts[lab.L], term.X, newVia.ID),
		}
		return s.takeProc(newProc)
	case step.SpawnSpec:
		newPE, err := s.Involve(PartSpec{Sig: term.Sig, Owner: proc.PID, TEs: term.CEs})
		if err != nil {
			return err
		}
		err = s.chnls.Transfer(id.Empty(), proc.PID, []chnl.ID{newPE.ID})
		if err != nil {
			s.log.Error("channel transfer failed",
				slog.Any("reason", err),
				slog.Any("from", id.Empty()),
				slog.Any("to", proc.PID),
				slog.Any("id", newPE.ID),
			)
			return err
		}
		s.log.Debug("transition taking succeeded")
		cfg.Add(newPE)
		cfg.Remove(term.CEs...)
		proc.Term = step.Subst(term.Cont, term.PE, newPE.ID)
		return s.takeProcWith(proc, cfg)
	case step.FwdSpec:
		viaID, ok := term.C.(chnl.ID)
		if !ok {
			err := chnl.ErrNotAnID(term.C)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curVia, ok := cfg.LookupCh(viaID)
		if !ok {
			err = chnl.ErrMissingInCfg(viaID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSt, ok := cfg.LookupSt(curVia.ID)
		if !ok {
			err = chnl.ErrMissingInCfg(curVia.ID)
			s.log.Error("transition taking failed",
				slog.Any("reason", err),
			)
			return err
		}
		curSem, err := s.steps.SelectByVID(curVia.ID)
		if err != nil {
			s.log.Error("step selection failed",
				slog.Any("reason", err),
				slog.Any("vid", curVia.ID),
			)
			return err
		}
		switch curSt.Pol() {
		case pol.Pos:
			switch sem := curSem.(type) {
			case step.SrvRoot:
				cID, ok := term.C.(chnl.ID)
				if !ok {
					err := chnl.ErrNotAnID(term.C)
					s.log.Error("transition taking failed", slog.Any("reason", err))
					return err
				}
				c, ok := cfg.LookupCh(cID)
				if !ok {
					err = chnl.ErrMissingInCfg(cID)
					s.log.Error("transition taking failed", slog.Any("reason", err))
					return err
				}
				err := s.chnls.Transfer(proc.PID, sem.PID, []chnl.ID{c.ID})
				if err != nil {
					s.log.Error("channel transfer failed",
						slog.Any("reason", err),
						slog.Any("from", proc.PID),
						slog.Any("to", sem.PID),
						slog.Any("id", c.ID),
					)
					return err
				}
				newProc := step.ProcRoot{
					ID:   id.New(),
					PID:  sem.PID,
					Term: step.Subst(sem.Cont, term.D, c.ID),
				}
				s.log.Debug("transition taking succeeded")
				return s.takeProc(newProc)
			case step.MsgRoot:
				dID, ok := term.C.(chnl.ID)
				if !ok {
					err := chnl.ErrNotAnID(term.D)
					s.log.Error("transition taking failed", slog.Any("reason", err))
					return err
				}
				d, ok := cfg.LookupCh(dID)
				if !ok {
					err = chnl.ErrMissingInCfg(dID)
					s.log.Error("transition taking failed", slog.Any("reason", err))
					return err
				}
				err := s.chnls.Transfer(proc.PID, sem.PID, []chnl.ID{d.ID})
				if err != nil {
					s.log.Error("channel transfer failed",
						slog.Any("reason", err),
						slog.Any("from", proc.PID),
						slog.Any("to", sem.PID),
						slog.Any("id", d.ID),
					)
					return err
				}
				newProc := step.ProcRoot{
					ID:   id.New(),
					PID:  sem.PID,
					Term: step.Subst(sem.Val, term.C, d.ID),
				}
				s.log.Debug("transition taking succeeded")
				return s.takeProc(newProc)
			case nil:
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
			default:
				panic(step.ErrRootTypeUnexpected(curSem))
			}
		case pol.Neg:
			switch sem := curSem.(type) {
			case step.SrvRoot:
				dID, ok := term.C.(chnl.ID)
				if !ok {
					err := chnl.ErrNotAnID(term.D)
					s.log.Error("transition taking failed", slog.Any("reason", err))
					return err
				}
				d, ok := cfg.LookupCh(dID)
				if !ok {
					err = chnl.ErrMissingInCfg(dID)
					s.log.Error("transition taking failed", slog.Any("reason", err))
					return err
				}
				err := s.chnls.Transfer(proc.PID, sem.PID, []chnl.ID{d.ID})
				if err != nil {
					s.log.Error("channel transfer failed",
						slog.Any("reason", err),
						slog.Any("from", proc.PID),
						slog.Any("to", sem.PID),
						slog.Any("id", d.ID),
					)
					return err
				}
				newProc := step.ProcRoot{
					ID:   id.New(),
					PID:  sem.PID,
					Term: step.Subst(sem.Cont, term.C, d.ID),
				}
				s.log.Debug("transition taking succeeded")
				return s.takeProc(newProc)
			case step.MsgRoot:
				cID, ok := term.C.(chnl.ID)
				if !ok {
					err := chnl.ErrNotAnID(term.C)
					s.log.Error("transition taking failed", slog.Any("reason", err))
					return err
				}
				c, ok := cfg.LookupCh(cID)
				if !ok {
					err = chnl.ErrMissingInCfg(cID)
					s.log.Error("transition taking failed", slog.Any("reason", err))
					return err
				}
				err := s.chnls.Transfer(proc.PID, sem.PID, []chnl.ID{c.ID})
				if err != nil {
					s.log.Error("channel transfer failed",
						slog.Any("reason", err),
						slog.Any("from", proc.PID),
						slog.Any("to", sem.PID),
						slog.Any("id", c.ID),
					)
					return err
				}
				newProc := step.ProcRoot{
					ID:   id.New(),
					PID:  sem.PID,
					Term: step.Subst(sem.Val, term.D, c.ID),
				}
				s.log.Debug("transition taking succeeded")
				return s.takeProc(newProc)
			case nil:
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

type repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectChildren(ID) ([]Ref, error)
	SelectSigs(ID) ([]sig.Ref, error)
}

// Kinship Relation
type KinshipSpec struct {
	ParentID ID
	ChildIDs []ID
}

type KinshipRoot struct {
	Parent   Ref
	Children []Ref
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// Participation aka lightweight Spawn
type PartSpec struct {
	Deal  id.ADT
	Sig   sig.ID
	Owner chnl.ID
	// Transferable Endpoints
	TEs []chnl.ID
}

// Transition
type TranSpec struct {
	Deal id.ADT
	// Process ID
	PID chnl.ID
	// Agent Access Key
	Key  ak.ADT
	Term step.Term
}

func (s *service) checkState(
	env Environment,
	ctx state.Context,
	pe state.EP,
	t step.Term,
) error {
	if pe.Z == t.Via() {
		return s.checkProvider(env, ctx, pe, t)
	}
	return s.checkClient(env, ctx, pe, t)
}

// aka checkExp
func (s *service) checkProvider(
	env Environment,
	ctx state.Context,
	pe state.EP,
	t step.Term,
) error {
	switch term := t.(type) {
	case step.CloseSpec:
		// check ctx
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
		wantSt, ok := pe.C.(state.TensorRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(pe.C, wantSt)
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
		err := state.CheckRoot(gotB, wantSt.B)
		if err != nil {
			return err
		}
		// no cont to check
		delete(ctx.Linear, term.B)
		pe.C = wantSt.C
		return nil
	case step.RecvSpec:
		// check via
		wantSt, ok := pe.C.(state.LolliRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(pe.C, wantSt)
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
		err := state.CheckRoot(gotY, wantSt.Y)
		if err != nil {
			return err
		}
		// check cont
		ctx.Linear[term.Y] = wantSt.Y
		pe.C = wantSt.Z
		return s.checkState(env, ctx, pe, term.Cont)
	case step.LabSpec:
		// check via
		wantSt, ok := pe.C.(state.PlusRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(pe.C, wantSt)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check label
		_, ok = wantSt.Choices[term.L]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %q", maps.Keys(wantSt.Choices), term.L)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// no cont to check
		return nil
	case step.CaseSpec:
		// check via
		wantSt, ok := pe.C.(state.WithRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(pe.C, wantSt)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check conts
		if len(term.Conts) != len(wantSt.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantSt.Choices), len(term.Conts))
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		for wantL, wantChoice := range wantSt.Choices {
			gotCont, ok := term.Conts[wantL]
			if !ok {
				err := fmt.Errorf("label mismatch: want %q, got nothing", wantL)
				s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
				return err
			}
			pe.C = wantChoice
			err := s.checkState(env, ctx, pe, gotCont)
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

func (s *service) checkClient(
	env Environment,
	ctx state.Context,
	pe state.EP,
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
		wantSt, ok := gotX.(state.OneRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotX, wantSt)
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
		wantSt, ok := gotA.(state.LolliRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotA, wantSt)
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
		err := state.CheckRoot(gotB, wantSt.Y)
		if err != nil {
			return err
		}
		// no cont to check
		delete(ctx.Linear, got.B)
		ctx.Linear[got.A] = wantSt.Z
		return nil
	case step.RecvSpec:
		// check via
		gotX, ok := ctx.Linear[got.X]
		if !ok {
			err := chnl.ErrMissingInCtx(got.X)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		wantSt, ok := gotX.(state.TensorRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotX, wantSt)
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
		err := state.CheckRoot(gotY, wantSt.B)
		if err != nil {
			return err
		}
		// check cont
		ctx.Linear[got.Y] = wantSt.B
		pe.C = wantSt.C
		return s.checkState(env, ctx, pe, got.Cont)
	case step.LabSpec:
		// check via
		gotA, ok := ctx.Linear[got.A]
		if !ok {
			err := chnl.ErrMissingInCtx(got.A)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		wantSt, ok := gotA.(state.WithRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotA, wantSt)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check label
		_, ok = wantSt.Choices[got.L]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %q", maps.Keys(wantSt.Choices), got.L)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// no cont to check
		return nil
	case step.CaseSpec:
		// check via
		gotX, ok := ctx.Linear[got.X]
		if !ok {
			err := chnl.ErrMissingInCtx(got.X)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		wantSt, ok := gotX.(state.PlusRoot)
		if !ok {
			err := state.ErrRootTypeMismatch(gotX, wantSt)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		// check conts
		if len(got.Conts) != len(wantSt.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantSt.Choices), len(got.Conts))
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		for wantL, wantChoice := range wantSt.Choices {
			gotCont, ok := got.Conts[wantL]
			if !ok {
				err := fmt.Errorf("label mismatch: want %q, got nothing", wantL)
				s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
				return err
			}
			ctx.Linear[got.X] = wantChoice
			err := s.checkState(env, ctx, pe, gotCont)
			if err != nil {
				return err
			}
		}
		return nil
	case step.SpawnSpec:
		if !env.Contains(got.Sig) {
			err := sig.ErrRootMissingInEnv(got.Sig)
			s.log.Error("type checking failed", slog.Any("reason", err), slog.Any("via", t.Via()))
			return err
		}
		wantCEs := env.LookupCEs(got.Sig)
		if len(got.CEs) != len(wantCEs) {
			err := fmt.Errorf("context mismatch: want %v items, got %v items", len(wantCEs), len(got.CEs))
			s.log.Error("type checking failed",
				slog.Any("reason", err),
				slog.Any("via", t.Via()),
				slog.Any("want", wantCEs),
				slog.Any("got", got.CEs),
			)
			return err
		}
		if len(got.CEs) == 0 {
			return nil
		}
		for i, gotCE := range got.CEs {
			gotSt := ctx.Linear[gotCE]
			err := state.CheckRoot(gotSt, wantCEs[i].C)
			if err != nil {
				s.log.Error("type checking failed",
					slog.Any("reason", err),
					slog.Any("via", t.Via()),
					slog.Any("want", wantCEs[i]),
					slog.Any("got", gotCE),
				)
				return err
			}
			delete(ctx.Linear, gotCE)
		}
		ctx.Linear[got.PE] = env.LookupPE(got.Sig).C
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
	linear := make(map[ph.ADT]state.Root, len(chnls))
	for _, ch := range chnls {
		linear[ch.ID] = states[*ch.StateID]
	}
	return state.Context{Linear: linear}
}
