package step

import (
	"fmt"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"

	"smecalculus/rolevod/app/seat"
)

type ID = id.ADT

type Ref interface {
	rID() ID
}

type ProcRef struct {
	ID ID
}

func (r ProcRef) rID() ID { return r.ID }

type MsgRef struct {
	ID ID
}

func (r MsgRef) rID() ID { return r.ID }

type SrvRef struct {
	ID ID
}

func (r SrvRef) rID() ID { return r.ID }

type Root interface {
	step()
}

// aka exec.Proc
type ProcRoot struct {
	ID   ID
	PID  chnl.ID
	Ctx  []chnl.ID
	Term Term
}

func (ProcRoot) step() {}

// aka exec.Msg
type MsgRoot struct {
	ID  ID
	PID chnl.ID
	VID chnl.ID
	Val Value
}

func (MsgRoot) step() {}

type SrvRoot struct {
	ID   ID
	PID  chnl.ID
	VID  chnl.ID
	Cont Continuation
}

func (SrvRoot) step() {}

type TbdRoot struct {
	ID  ID
	PID chnl.ID
	VID chnl.ID
	Act Action
}

func (TbdRoot) step() {}

// aka Expression
type Term interface {
	Via() core.Placeholder
}

// aka ast.Msg
type Value interface {
	Term
	val()
}

type Continuation interface {
	Term
	cont()
}

type Action interface {
	Term
	act()
}

type CloseSpec struct {
	A core.Placeholder
}

func (CloseSpec) val() {}

func (s CloseSpec) Via() core.Placeholder { return s.A }

type WaitSpec struct {
	X    core.Placeholder
	Cont Term
}

func (WaitSpec) cont() {}

func (s WaitSpec) Via() core.Placeholder { return s.X }

type SendSpec struct {
	A core.Placeholder // via
	B core.Placeholder // value
	// Cont  Term
}

func (SendSpec) val() {}

func (s SendSpec) Via() core.Placeholder { return s.A }

type RecvSpec struct {
	X    core.Placeholder // via
	Y    core.Placeholder // value
	Cont Term
}

func (RecvSpec) cont() {}

func (s RecvSpec) Via() core.Placeholder { return s.X }

type LabSpec struct {
	C core.Placeholder
	L core.Label
	// Cont Term
}

func (LabSpec) val() {}

func (s LabSpec) Via() core.Placeholder { return s.C }

type CaseSpec struct {
	Z     core.Placeholder
	Conts map[core.Label]Term
}

func (CaseSpec) cont() {}

func (s CaseSpec) Via() core.Placeholder { return s.Z }

type CTASpec struct {
	AK  ak.ADT
	SID seat.ID
	Ctx map[chnl.Name]chnl.ID
}

func (s CTASpec) act() {}

func (s CTASpec) Via() core.Placeholder { return s.SID }

// aka ExpName
type RecurSpec struct {
	C    chnl.ID
	Ctx  []chnl.ID
	Seat sym.ADT
}

func (s RecurSpec) Via() core.Placeholder { return s.C }

type FwdSpec struct {
	C core.Placeholder // from
	D core.Placeholder // to
}

func (FwdSpec) val() {}

func (FwdSpec) cont() {}

func (s FwdSpec) Via() core.Placeholder { return s.C }

type SpawnSpec struct {
	Z      core.Placeholder
	Ctx    []chnl.ID
	Cont   Term
	SeatID id.ADT
}

func (s SpawnSpec) Via() core.Placeholder { return s.Z }

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectByPID(chnl.ID) (Root, error)
	SelectByVID(chnl.ID) (Root, error)
}

func CollectEnv(t Term) []id.ADT {
	return collectEnvRec(t, []id.ADT{})
}

func collectEnvRec(t Term, env []id.ADT) []id.ADT {
	switch term := t.(type) {
	case RecvSpec:
		return collectEnvRec(term.Cont, env)
	case CaseSpec:
		for _, cont := range term.Conts {
			env = collectEnvRec(cont, env)
		}
		return env
	case SpawnSpec:
		return collectEnvRec(term.Cont, append(env, term.SeatID))
	default:
		return env
	}
}

func CollectCtx(pid chnl.ID, t Term) []chnl.ID {
	return collectCtxRec(pid, t, nil)
}

func collectCtxRec(pid chnl.ID, t Term, ctx []chnl.ID) []chnl.ID {
	switch term := t.(type) {
	case WaitSpec:
		x, ok := term.X.(chnl.ID)
		if ok && x != pid {
			ctx = append(ctx, x)
		}
		return collectCtxRec(pid, term.Cont, ctx)
	case SendSpec:
		a, ok := term.A.(chnl.ID)
		if ok && a != pid {
			ctx = append(ctx, a)
		}
		b, ok := term.B.(chnl.ID)
		if ok {
			ctx = append(ctx, b)
		}
		return ctx
	case RecvSpec:
		x, ok := term.X.(chnl.ID)
		if ok && x != pid {
			ctx = append(ctx, x)
		}
		y, ok := term.Y.(chnl.ID)
		if ok {
			ctx = append(ctx, y)
		}
		return collectCtxRec(pid, term.Cont, ctx)
	case LabSpec:
		c, ok := term.C.(chnl.ID)
		if ok && c != pid {
			ctx = append(ctx, c)
		}
		return ctx
	case CaseSpec:
		z, ok := term.Z.(chnl.ID)
		if ok && z != pid {
			ctx = append(ctx, z)
		}
		for _, cont := range term.Conts {
			ctx = collectCtxRec(pid, cont, ctx)
		}
		return ctx
	case FwdSpec:
		d, ok := term.D.(chnl.ID)
		if ok {
			ctx = append(ctx, d)
		}
		return ctx
	case SpawnSpec:
		return collectCtxRec(pid, term.Cont, append(ctx, term.Ctx...))
	default:
		return ctx
	}
}

func Subst(t Term, ph core.Placeholder, val chnl.ID) Term {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case CloseSpec:
		if ph == term.A {
			term.A = val
		}
		return term
	case WaitSpec:
		if ph == term.X {
			term.X = val
		}
		term.Cont = Subst(term.Cont, ph, val)
		return term
	case SendSpec:
		if ph == term.A {
			term.A = val
		}
		if ph == term.B {
			term.B = val
		}
		return term
	default:
		panic(ErrTermTypeUnexpected(t))
	}
}

func ErrDoesNotExist(want ID) error {
	return fmt.Errorf("root doesn't exist: %v", want)
}

func ErrRootTypeUnexpected(got Root) error {
	return fmt.Errorf("root type unexpected: %T", got)
}

func ErrRootTypeMismatch(got, want Root) error {
	return fmt.Errorf("root type mismatch: want %T, got %T", want, got)
}

func ErrTermTypeUnexpected(got Term) error {
	return fmt.Errorf("term type unexpected: %T", got)
}

func ErrTermTypeMismatch(got, want Term) error {
	return fmt.Errorf("term type mismatch: want %T, got %T", want, got)
}

func ErrTermValueNil(pid chnl.ID) error {
	return fmt.Errorf("proc %q term is nil", pid)
}

func ErrValTypeUnexpected(got Value) error {
	return fmt.Errorf("value type unexpected: %T", got)
}

func ErrContTypeUnexpected(got Continuation) error {
	return fmt.Errorf("continuation type unexpected: %T", got)
}
