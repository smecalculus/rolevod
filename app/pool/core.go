package pool

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/rev"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/proc"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/sig"
)

type ID = id.ADT
type Rev = rev.ADT
type Title = string

type Spec struct {
	Title  string
	SupID  id.ADT
	DepIDs []sig.ID
}

type Ref struct {
	PoolID id.ADT
	Title  string
}

type Root struct {
	PoolID id.ADT
	Title  string
	SupID  id.ADT
	Revs   []rev.ADT
}

const (
	rootRev rev.Knd = iota
	procRev
)

type RootMod struct {
	PoolID id.ADT
	Rev    rev.ADT
	K      rev.Knd
}

type SubSnap struct {
	PoolID id.ADT
	Title  string
	Subs   []Ref
}

type AssetSnap struct {
	PoolID id.ADT
	Title  string
	EPs    []proc.EP
}

type AssetMod struct {
	OutPoolID id.ADT
	InPoolID  id.ADT
	Rev       rev.ADT
	EPs       []proc.EP
}

type LiabSnap struct {
	PoolID id.ADT
	Title  string
	EP     proc.EP
}

type LiabMod struct {
	OutPoolID id.ADT
	InPoolID  id.ADT
	Rev       rev.ADT
	EP        proc.EP
}

type TranSpec struct {
	PoolID id.ADT
	ProcID id.ADT
	Term   step.Term
}

type Environment struct {
	sigs   map[sig.ID]sig.Root
	roles  map[role.QN]role.Root
	states map[state.ID]state.Root
}

func (e Environment) Contains(id sig.ID) bool {
	_, ok := e.sigs[id]
	return ok
}

func (e Environment) LookupPE(id sig.ID) state.EP {
	decl := e.sigs[id]
	role := e.roles[decl.PE.Link]
	return state.EP{Z: decl.PE.Link, C: e.states[role.StateID]}
}

func (e Environment) LookupCEs(id sig.ID) []state.EP {
	decl := e.sigs[id]
	ces := []state.EP{}
	for _, ce := range decl.CEs {
		role := e.roles[decl.PE.Link]
		ces = append(ces, state.EP{Z: ce.Link, C: e.states[role.StateID]})
	}
	return ces
}

// Port
type API interface {
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (SubSnap, error)
	RetreiveRefs() ([]Ref, error)
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

type service struct {
	pools    Repo
	sigs     sig.Repo
	roles    role.Repo
	states   state.Repo
	operator data.Operator
	log      *slog.Logger
}

func newService(
	pools Repo,
	sigs sig.Repo,
	roles role.Repo,
	states state.Repo,
	operator data.Operator,
	l *slog.Logger,
) *service {
	name := slog.String("name", "poolService")
	return &service{pools, sigs, roles, states, operator, l.With(name)}
}

func (s *service) Create(spec Spec) (_ Root, err error) {
	ctx := context.Background()
	s.log.Debug("creation started", slog.Any("spec", spec))
	root := Root{
		PoolID: id.New(),
		Revs:   []rev.ADT{rev.Initial()},
		Title:  spec.Title,
		SupID:  spec.SupID,
	}
	s.operator.Explicit(ctx, func(ds data.Source) error {
		err = s.pools.Insert(ds, root)
		return err
	})
	if err != nil {
		s.log.Error("creation failed")
		return Root{}, err
	}
	s.log.Debug("creation succeeded", slog.Any("id", root.PoolID))
	return root, nil
}

func (s *service) Take(spec TranSpec) (err error) {
	ctx := context.Background()
	// initial iteration
	curProcID := spec.ProcID
	curTerm := spec.Term
	for curTerm != nil {
		var procSnap proc.Snap
		s.operator.Implicit(ctx, func(ds data.Source) {
			procSnap, err = s.pools.SelectProc(ds, curProcID)
		})
		idAttr := slog.Any("procID", curProcID)
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		if len(procSnap.EPs) == 0 {
			s.log.Error("taking failed", idAttr)
			return err
		}
		sigIDs := step.CollectEnv(curTerm)
		var sigs map[sig.ID]sig.Root
		s.operator.Implicit(ctx, func(ds data.Source) {
			sigs, err = s.sigs.SelectEnv(ds, sigIDs)
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("sigs", sigIDs))
			return err
		}
		roleQNs := sig.CollectEnv(maps.Values(sigs))
		var roles map[role.QN]role.Root
		s.operator.Implicit(ctx, func(ds data.Source) {
			roles, err = s.roles.SelectEnv(ds, roleQNs)
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("roles", roleQNs))
			return err
		}
		envIDs := role.CollectEnv(maps.Values(roles))
		ctxIDs := CollectCtx(maps.Values(procSnap.EPs))
		var states map[state.ID]state.Root
		s.operator.Implicit(ctx, func(ds data.Source) {
			states, err = s.states.SelectEnv(ds, append(envIDs, ctxIDs...))
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("env", envIDs), slog.Any("ctx", ctxIDs))
			return err
		}
		environ := Environment{sigs, roles, states}
		context := convertToCtx(maps.Values(procSnap.EPs), states)
		// type checking
		err = s.checkState(spec.PoolID, environ, context, procSnap, curTerm)
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		// step taking
		nextTerm, procMod, err := s.takeWith(context, procSnap, curTerm)
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		s.operator.Explicit(ctx, func(ds data.Source) error {
			err := s.pools.UpdateProc(ds, procMod)
			if err != nil {
				s.log.Error("taking failed", idAttr)
				return err
			}
			return nil
		})
		// next iteration
		curProcID = id.Nil
		curTerm = nextTerm
	}
	return nil
}

func (s *service) takeWith(
	ctx state.Context,
	snap proc.Snap,
	t step.Term,
) (
	_ step.Term,
	mod proc.Mod,
	_ error,
) {
	switch term := t.(type) {
	case step.CloseSpec:
		viaEP, ok := snap.EPs[term.A]
		if !ok {
			err := chnl.ErrMissingInCfg(term.A)
			s.log.Error("taking failed")
			return nil, mod, err
		}
		epStep, ok := snap.Steps[viaEP.ChnlID]
		viaAttr := slog.Any("via", viaEP.ChnlID)
		if !ok {
			err := step.ErrMissingInCfg(viaEP.ChnlID)
			s.log.Error("taking failed", viaAttr)
			return nil, mod, err
		}
		if epStep == nil {
			newStep := step.MsgRoot2{
				ProcID: viaEP.ProcID,
				ChnlID: viaEP.ChnlID,
				Val:    term,
			}
			mod.Steps = append(mod.Steps, newStep)
			s.log.Debug("transition taking half done", slog.Any("msg", newStep))
			return nil, mod, nil
		}
		srvStep, ok := epStep.(step.SrvRoot2)
		if !ok {
			err := step.ErrRootTypeUnexpected(epStep)
			s.log.Error("taking failed")
			return nil, mod, err
		}
		waitSpec, ok := srvStep.Cont.(step.WaitSpec)
		if !ok {
			err := fmt.Errorf("cont type mismatch: want %T, got %T", waitSpec, srvStep.Cont)
			s.log.Error("taking failed", slog.Any("cont", srvStep.Cont))
			return nil, mod, err
		}
		// close channel
		newBnd := proc.Binding{
			ProcID:  viaEP.ProcID,
			ProcPH:  viaEP.ChnlPH,
			ChnlID:  id.Nil,
			StateID: id.Nil,
			Rev:     viaEP.PrvdRevs[1] + 1,
		}
		mod.Bnds = append(mod.Bnds, newBnd)
		mod.PoolID = viaEP.PrvdID
		mod.Rev = viaEP.PrvdRevs[procRev]
		s.log.Debug("transition taking succeeded")
		return waitSpec.Cont, mod, nil
	default:
		panic(step.ErrTermTypeUnexpected(t))
	}
}

func (s *service) Retrieve(poolID id.ADT) (snap SubSnap, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) {
		snap, err = s.pools.SelectSubs(ds, poolID)
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("id", poolID))
		return SubSnap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveRefs() (refs []Ref, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) {
		refs, err = s.pools.SelectRefs(ds)
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

func CollectCtx(eps []proc.EP) []state.ID {
	return nil
}

func convertToCtx(eps []proc.EP, states map[state.ID]state.Root) state.Context {
	linear := make(map[ph.ADT]state.Root, len(eps))
	for _, ep := range eps {
		linear[ep.ChnlPH] = states[ep.StateID]
	}
	return state.Context{Linear: linear}
}

func convertToCfg(eps []proc.EP) map[ph.ADT]proc.EP {
	cfg := make(map[ph.ADT]proc.EP, len(eps))
	for _, ep := range eps {
		cfg[ep.ChnlPH] = ep
	}
	return cfg
}

func (s *service) checkState(
	poolID id.ADT,
	env Environment,
	ctx state.Context,
	acl proc.Snap,
	t step.Term,
) error {
	ep, ok := acl.EPs[t.Via()]
	if !ok {
		panic("no via in acl")
	}
	if ep.PrvdID == ep.ClntID {
		panic("can not be equal")
	}
	switch poolID {
	case ep.PrvdID:
		return s.checkProvider(poolID, env, ctx, acl, t)
	case ep.ClntID:
		return s.checkClient(poolID, env, ctx, acl, t)
	default:
		s.log.Error("state checking failed", slog.Any("id", poolID))
		return fmt.Errorf("unknown pool id")
	}
}

func (s *service) checkProvider(
	poolID id.ADT,
	env Environment,
	ctx state.Context,
	acl proc.Snap,
	t step.Term,
) error {
	return nil
}

func (s *service) checkClient(
	poolID id.ADT,
	env Environment,
	ctx state.Context,
	acl proc.Snap,
	t step.Term,
) error {
	return nil
}

// Port
type Repo interface {
	Insert(data.Source, Root) error
	SelectRefs(data.Source) ([]Ref, error)
	SelectSubs(data.Source, id.ADT) (SubSnap, error)
	SelectAssets(data.Source, id.ADT) (AssetSnap, error)
	SelectEPsByProcID(data.Source, id.ADT) ([]proc.EP, error)
	SelectProc(data.Source, id.ADT) (proc.Snap, error)
	UpdateProc(data.Source, proc.Mod) error
	UpdateAssets(data.Source, AssetMod) error
	Transfer(source data.Source, giver id.ADT, taker id.ADT, pids []chnl.ID) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)

func errOptimisticUpdate(got rev.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
