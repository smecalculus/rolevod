package step

import (
	"fmt"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

var (
	// TODO подобрать другой юнит
	Unit = Ref{}
)

type ID interface{}

type Ref struct {
	ID id.ADT[ID]
}

func (Ref) val() {}

type root interface {
	step()
}

// aka exec.Proc
type Process struct {
	ID    id.ADT[ID]
	PreID id.ADT[ID]
	Term  Term
}

func (Process) step() {}

// aka exec.Msg
type Message struct {
	ID    id.ADT[ID]
	PreID id.ADT[ID]
	ViaID id.ADT[chnl.ID]
	Val   Value
}

func (Message) step() {}

type Service struct {
	ID    id.ADT[ID]
	PreID id.ADT[ID]
	ViaID id.ADT[chnl.ID]
	Cont  Continuation
}

func (Service) step() {}

type Label string

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

type Fwd struct {
	A chnl.Ref // from
	C chnl.Ref // to
}

func (Fwd) term() {}

type Spawn struct {
	Name string
	C    chnl.Var
	Ctx  []chnl.Ref
	Cont Term
}

func (Spawn) term() {}

type Lab struct {
	C chnl.Ref
	L Label
	// Cont Term
}

func (Lab) term() {}
func (Lab) val()  {}

type Case struct {
	X     chnl.Ref
	Conts map[Label]Term
}

func (Case) term() {}
func (Case) cont() {}

type Send struct {
	A chnl.Ref // channel
	B chnl.Ref // value
	// Cont  Term
}

func (Send) term() {}
func (Send) val()  {}

type Recv struct {
	X    chnl.Ref // channel
	Y    chnl.Ref // value
	Cont Term
}

func (Recv) term() {}
func (Recv) cont() {}

type Close struct {
	A chnl.Ref
}

func (Close) term() {}
func (Close) val()  {}

type Wait struct {
	X    chnl.Ref
	Cont Term
}

func (Wait) term() {}
func (Wait) cont() {}

// aka ExpName
type ExpRef struct {
	Name string
	C    chnl.Var
	Ctx  []chnl.Var
}

func (ExpRef) term() {}

type Repo[T root] interface {
	Insert(root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT[ID]) (*T, error)
	SelectByChID(id.ADT[chnl.ID]) (*T, error)
}

func Subst(t Term, from chnl.Ref, to chnl.Ref) {
	switch term := t.(type) {
	case Close:
		if term.A != from {
			return
		}
		term.A = to
	default:
		panic(ErrUnexpectedTerm(term))
	}
}

func ErrUnexpectedStep(s root) error {
	return fmt.Errorf("unexpected steo %#v", s)
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

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
