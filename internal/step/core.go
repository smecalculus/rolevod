package step

import (
	"fmt"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

type ID interface{}

type Ref interface {
	rootID() id.ADT[ID]
}

type ProcRef id.ADT[ID]

func (r ProcRef) rootID() id.ADT[ID] { return id.ADT[ID](r) }

type MsgRef id.ADT[ID]

func (r MsgRef) rootID() id.ADT[ID] { return id.ADT[ID](r) }

type SrvRef id.ADT[ID]

func (r SrvRef) rootID() id.ADT[ID] { return id.ADT[ID](r) }

type root interface {
	step()
}

// aka exec.Proc
type ProcRoot struct {
	ID    id.ADT[ID]
	PreID id.ADT[ID]
	Term  Term
}

func (ProcRoot) step() {}

// aka exec.Msg
type MsgRoot struct {
	ID    id.ADT[ID]
	PreID id.ADT[ID]
	ViaID id.ADT[chnl.ID]
	Val   Value
}

func (MsgRoot) step() {}

type SrvRoot struct {
	ID    id.ADT[ID]
	PreID id.ADT[ID]
	ViaID id.ADT[chnl.ID]
	Cont  Continuation
}

func (SrvRoot) step() {}

type Label string

// aka Expression
type Term interface {
	term()
	// subst(variable chnl.Ref, value chnl.Ref)
}

// aka ast.Msg
type Value interface {
	val()
}

type Continuation interface {
	cont()
}

type FwdSpec struct {
	A chnl.Ref // from
	C chnl.Ref // to
}

func (FwdSpec) term() {}

// func (s *FwdSpec) subst(variable chnl.Ref, value chnl.Ref) {
// }

type SpawnSpec struct {
	Name string
	C    chnl.Var
	Ctx  []chnl.Ref
	Cont Term
}

func (SpawnSpec) term() {}

// func (s *SpawnSpec) subst(variable chnl.Ref, value chnl.Ref) {
// }

type LabSpec struct {
	C chnl.Ref
	L Label
	// Cont Term
}

func (LabSpec) term() {}
func (LabSpec) val()  {}

// func (s *LabSpec) subst(variable chnl.Ref, value chnl.Ref) {
// }

type CaseSpec struct {
	X     chnl.Ref
	Conts map[Label]Term
}

func (CaseSpec) term() {}
func (CaseSpec) cont() {}

// func (s *CaseSpec) subst(variable chnl.Ref, value chnl.Ref) {
// }

type SendSpec struct {
	A chnl.Ref // channel
	B chnl.Ref // value
	// Cont  Term
}

func (SendSpec) term() {}
func (SendSpec) val()  {}

// func (s *SendSpec) subst(variable chnl.Ref, value chnl.Ref) {
// }

type RecvSpec struct {
	X    chnl.Ref // channel
	Y    chnl.Ref // value
	Cont Term
}

func (RecvSpec) term() {}
func (RecvSpec) cont() {}

// func (s *RecvSpec) subst(variable chnl.Ref, value chnl.Ref) {
// }

type CloseSpec struct {
	A chnl.Ref
}

func (CloseSpec) term() {}
func (CloseSpec) val()  {}

// func (s *CloseSpec) subst(variable chnl.Ref, value chnl.Ref) {
// }

type WaitSpec struct {
	X    chnl.Ref
	Cont Term
}

func (WaitSpec) term() {}
func (WaitSpec) cont() {}

// func (s *WaitSpec) subst(variable chnl.Ref, value chnl.Ref) {
// 	if s.X == variable {
// 		s.X = value
// 	}
// }

// aka ExpName
type RecSpec struct {
	Name string
	C    chnl.Var
	Ctx  []chnl.Var
}

func (RecSpec) term() {}

// func (s *RecSpec) subst(variable chnl.Ref, value chnl.Ref) {
// }

type Repo[T root] interface {
	Insert(root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT[ID]) (*T, error)
	SelectByCh(id.ADT[chnl.ID]) (*T, error)
}

func Subst(t Term, varbl chnl.Ref, value chnl.Ref) {
	if t == nil {
		return
	}
	switch term := t.(type) {
	case CloseSpec:
		if term.A == varbl {
			term.A = value
		}
	case WaitSpec:
		if term.X == varbl {
			term.X = value
		}
		Subst(term.Cont, varbl, value)
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

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
