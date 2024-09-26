package step

import (
	"fmt"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
)

type ID interface{}

type Ref interface {
	rootID() id.ADT
}

type ProcRef id.ADT

func (r ProcRef) rootID() id.ADT { return id.ADT(r) }

type MsgRef id.ADT

func (r MsgRef) rootID() id.ADT { return id.ADT(r) }

type SrvRef id.ADT

func (r SrvRef) rootID() id.ADT { return id.ADT(r) }

type root interface {
	step()
}

// aka exec.Proc
type ProcRoot struct {
	ID   id.ADT
	Term Term
}

func (ProcRoot) step() {}

// aka exec.Msg
type MsgRoot struct {
	ID    id.ADT
	PreID id.ADT
	ViaID id.ADT
	Val   Value
}

func (MsgRoot) step() {}

type SrvRoot struct {
	ID    id.ADT
	PreID id.ADT
	ViaID id.ADT
	Cont  Continuation
}

func (SrvRoot) step() {}

// aka Expression
type Term interface {
	term()
}

// aka ast.Msg
type Value interface {
	val()
}

type Continuation interface {
	cont()
}

type Substitutable interface {
	Subst(chnl.Ref, chnl.Ref)
}

type FwdSpec struct {
	A chnl.Ref // from
	C chnl.Ref // to
}

func (FwdSpec) term() {}

type SpawnSpec struct {
	Name string
	C    chnl.Sym
	Ctx  []chnl.Ref
	Cont Term
}

func (SpawnSpec) term() {}

type LabSpec struct {
	C chnl.Ref
	L state.Label
	// Cont Term
}

func (LabSpec) term() {}
func (LabSpec) val()  {}

type CaseSpec struct {
	X        chnl.Ref
	Branches map[state.Label]Term
}

func (CaseSpec) term() {}
func (CaseSpec) cont() {}

type SendSpec struct {
	A chnl.Ref // channel
	B chnl.Ref // value
	// Cont  Term
}

func (SendSpec) term() {}
func (SendSpec) val()  {}

type RecvSpec struct {
	X    chnl.Ref // channel
	Y    chnl.Ref // value
	Cont Term
}

func (RecvSpec) term() {}
func (RecvSpec) cont() {}

type CloseSpec struct {
	A chnl.Ref
}

func (CloseSpec) term() {}
func (CloseSpec) val()  {}

func (t *CloseSpec) Subst(varRef chnl.Ref, valRef chnl.Ref) {
	if varRef == t.A {
		t.A = valRef
	}
}

type WaitSpec struct {
	X    chnl.Ref
	Cont Term
}

func (WaitSpec) term() {}
func (WaitSpec) cont() {}

// aka ExpName
type RecurSpec struct {
	Name string
	C    chnl.Sym
	Ctx  []chnl.Sym
}

func (RecurSpec) term() {}

type Repo[T root] interface {
	Insert(root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT) (*T, error)
	SelectByCh(id.ADT) (*T, error)
}

func Subst(t Term, varRef chnl.Ref, valRef chnl.Ref) Term {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case CloseSpec:
		if varRef == term.A {
			term.A = valRef
		}
		return term
	case WaitSpec:
		if varRef == term.X {
			term.X = valRef
		}
		term.Cont = Subst(term.Cont, varRef, valRef)
		return term
	default:
		panic(ErrUnexpectedTerm(t))
	}
}

func ErrUnexpectedStep(s root) error {
	return fmt.Errorf("unexpected step %#v", s)
}

func ErrUnexpectedTerm(t Term) error {
	return fmt.Errorf("unexpected term %#v", t)
}

func ErrUnexpectedValue(v Value) error {
	return fmt.Errorf("unexpected value %#v", v)
}

func ErrUnexpectedCont(c Continuation) error {
	return fmt.Errorf("unexpected continuation %#v", c)
}
