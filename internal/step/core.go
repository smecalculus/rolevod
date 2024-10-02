package step

import (
	"fmt"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
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
	ID  ID
	VID chnl.ID
	Val Value
}

func (MsgRoot) step() {}

type SrvRoot struct {
	ID   ID
	VID  chnl.ID
	Cont Continuation
}

func (SrvRoot) step() {}

// aka Expression
type Term interface {
	term()
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

type FwdSpec struct {
	A chnl.ID // from
	C chnl.ID // to
}

func (FwdSpec) term() {}

type SpawnSpec struct {
	DecID id.ADT // seat id
	C     chnl.ID
	Ctx   map[chnl.Name]chnl.ID
	Cont  Term
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
	Z     chnl.ID
	Conts map[state.Label]Term
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

type WaitSpec struct {
	X    chnl.ID
	Cont Term
}

func (WaitSpec) term() {}
func (WaitSpec) cont() {}

// aka ExpName
type RecurSpec struct {
	Name string
	C    chnl.Name
	Ctx  []chnl.Name
}

func (RecurSpec) term() {}

type Repo[T root] interface {
	Insert(root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (*T, error)
	SelectByCh(chnl.ID) (*T, error)
}

func CollectChnlIDs(t Term, ids []chnl.ID) []chnl.ID {
	switch term := t.(type) {
	case CloseSpec:
		return append(ids, term.A)
	case WaitSpec:
		return CollectChnlIDs(term.Cont, append(ids, term.X))
	case SendSpec:
		return append(ids, term.A, term.B)
	case RecvSpec:
		return CollectChnlIDs(term.Cont, append(ids, term.X, term.Y))
	default:
		panic(ErrUnexpectedTerm(t))
	}
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
