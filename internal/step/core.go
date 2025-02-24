package step

import (
	"fmt"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
)

type ID = id.ADT

type Ref interface {
	id.Identifiable
}

type ProcRef struct {
	ID ID
}

func (r ProcRef) Ident() ID { return r.ID }

type MsgRef struct {
	ID ID
}

func (r MsgRef) Ident() ID { return r.ID }

type SrvRef struct {
	ID ID
}

func (r SrvRef) Ident() ID { return r.ID }

type Root interface {
	step()
}

// aka exec.Proc
type ProcRoot struct {
	ID   ID
	PID  chnl.ID
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
	Via() ph.ADT
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
	A ph.ADT
}

func (CloseSpec) val() {}

func (s CloseSpec) Via() ph.ADT { return s.A }

type WaitSpec struct {
	X    ph.ADT
	Cont Term
}

func (WaitSpec) cont() {}

func (s WaitSpec) Via() ph.ADT { return s.X }

type SendSpec struct {
	A ph.ADT // via
	B ph.ADT // value
	// Cont  Term
}

func (SendSpec) val() {}

func (s SendSpec) Via() ph.ADT { return s.A }

type RecvSpec struct {
	X    ph.ADT // via
	Y    ph.ADT // value
	Cont Term
}

func (RecvSpec) cont() {}

func (s RecvSpec) Via() ph.ADT { return s.X }

type LabSpec struct {
	A ph.ADT
	L core.Label
	// Cont Term
}

func (LabSpec) val() {}

func (s LabSpec) Via() ph.ADT { return s.A }

type CaseSpec struct {
	X     ph.ADT
	Conts map[core.Label]Term
}

func (CaseSpec) cont() {}

func (s CaseSpec) Via() ph.ADT { return s.X }

type CTASpec struct {
	AK  ak.ADT
	Sig id.ADT
}

func (s CTASpec) act() {}

func (s CTASpec) Via() ph.ADT { return s.Sig }

// aka ExpName
type LinkSpec struct {
	PE  chnl.ID
	CEs []chnl.ID
	Sig sym.ADT
}

func (s LinkSpec) Via() ph.ADT { return s.PE }

type FwdSpec struct {
	C ph.ADT // from
	D ph.ADT // to
}

func (FwdSpec) val() {}

func (FwdSpec) cont() {}

func (s FwdSpec) Via() ph.ADT { return s.C }

type SpawnSpec struct {
	PE   ph.ADT
	CEs  []chnl.ID
	Cont Term
	Sig  id.ADT
}

func (s SpawnSpec) Via() ph.ADT { return s.PE }

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
		return collectEnvRec(term.Cont, append(env, term.Sig))
	default:
		return env
	}
}

func CollectCtx(pe chnl.ID, t Term) []chnl.ID {
	return collectCEsRec(pe, t, nil)
}

func collectCEsRec(pe chnl.ID, t Term, ces []chnl.ID) []chnl.ID {
	switch term := t.(type) {
	case WaitSpec:
		x, ok := term.X.(chnl.ID)
		if ok && x != pe {
			ces = append(ces, x)
		}
		return collectCEsRec(pe, term.Cont, ces)
	case SendSpec:
		a, ok := term.A.(chnl.ID)
		if ok && a != pe {
			ces = append(ces, a)
		}
		b, ok := term.B.(chnl.ID)
		if ok {
			ces = append(ces, b)
		}
		return ces
	case RecvSpec:
		x, ok := term.X.(chnl.ID)
		if ok && x != pe {
			ces = append(ces, x)
		}
		y, ok := term.Y.(chnl.ID)
		if ok {
			ces = append(ces, y)
		}
		return collectCEsRec(pe, term.Cont, ces)
	case LabSpec:
		a, ok := term.A.(chnl.ID)
		if ok && a != pe {
			ces = append(ces, a)
		}
		return ces
	case CaseSpec:
		x, ok := term.X.(chnl.ID)
		if ok && x != pe {
			ces = append(ces, x)
		}
		for _, cont := range term.Conts {
			ces = collectCEsRec(pe, cont, ces)
		}
		return ces
	case FwdSpec:
		d, ok := term.D.(chnl.ID)
		if ok {
			ces = append(ces, d)
		}
		return ces
	case SpawnSpec:
		return collectCEsRec(pe, term.Cont, append(ces, term.CEs...))
	default:
		return ces
	}
}

func Subst(t Term, ph ph.ADT, val chnl.ID) Term {
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
