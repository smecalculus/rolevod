package step

import (
	"fmt"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
)

type ID = id.ADT

type Ref interface {
	rootID() ID
}

type ProcRef ID

func (r ProcRef) rootID() ID { return ID(r) }

type MsgRef ID

func (r MsgRef) rootID() ID { return ID(r) }

type SrvRef ID

func (r SrvRef) rootID() ID { return ID(r) }

type root interface {
	step()
}

// aka exec.Proc
type ProcRoot struct {
	ID   ID
	Term Term
}

func (ProcRoot) step() {}

// aka exec.Msg
type MsgRoot struct {
	ID    ID
	ViaID chnl.ID
	Val   Value
}

func (MsgRoot) step() {}

type SrvRoot struct {
	ID    ID
	ViaID chnl.ID
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
	Subst(chnl.ID, chnl.ID)
}

type FwdSpec struct {
	A chnl.ID // from
	C chnl.ID // to
}

func (FwdSpec) term() {}

type SpawnSpec struct {
	Name string
	C    chnl.Sym
	Ctx  []chnl.ID
	Cont Term
}

func (SpawnSpec) term() {}

type LabSpec struct {
	C chnl.ID
	L state.Label
	// Cont Term
}

func (LabSpec) term() {}
func (LabSpec) val()  {}

type CaseSpec struct {
	X        chnl.ID
	Branches map[state.Label]Term
}

func (CaseSpec) term() {}
func (CaseSpec) cont() {}

type SendSpec struct {
	A chnl.ID // channel
	B chnl.ID // value
	// Cont  Term
}

func (SendSpec) term() {}
func (SendSpec) val()  {}

type RecvSpec struct {
	X    chnl.ID // channel
	Y    chnl.ID // value
	Cont Term
}

func (RecvSpec) term() {}
func (RecvSpec) cont() {}

type CloseSpec struct {
	A chnl.ID
}

func (CloseSpec) term() {}
func (CloseSpec) val()  {}

func (t *CloseSpec) Subst(varRef chnl.ID, valRef chnl.ID) {
	if varRef == t.A {
		t.A = valRef
	}
}

type WaitSpec struct {
	X    chnl.ID
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
	SelectByID(ID) (*T, error)
	SelectByCh(ID) (*T, error)
}

func Subst(t Term, varID chnl.ID, valID chnl.ID) Term {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case CloseSpec:
		if varID == term.A {
			term.A = valID
		}
		return term
	case WaitSpec:
		if varID == term.X {
			term.X = valID
		}
		term.Cont = Subst(term.Cont, varID, valID)
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
