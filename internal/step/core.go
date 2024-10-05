package step

import (
	"fmt"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

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

type Root interface {
	step()
}

// aka exec.Proc
type ProcRoot struct {
	ID   ID
	PID  chnl.ID
	Ctx  map[chnl.Name]chnl.ID
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
	PID  chnl.ID
	VID  chnl.ID
	Cont Continuation
}

func (SrvRoot) step() {}

// aka Expression
type Term interface {
	VID() chnl.ID
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

type LabSpec struct {
	C chnl.ID
	L state.Label
	// Cont Term
}

func (LabSpec) val() {}

func (s LabSpec) VID() chnl.ID          { return s.C }
func (s LabSpec) Via() core.Placeholder { return s.C }

type CaseSpec struct {
	Z     chnl.ID
	Conts map[state.Label]Term
}

func (CaseSpec) cont() {}

func (s CaseSpec) VID() chnl.ID          { return s.Z }
func (s CaseSpec) Via() core.Placeholder { return s.Z }

type SendSpec struct {
	A chnl.ID // channel
	B chnl.ID // value
	// Cont  Term
}

func (SendSpec) val() {}

func (s SendSpec) VID() chnl.ID          { return s.A }
func (s SendSpec) Via() core.Placeholder { return s.A }

type RecvSpec struct {
	X    chnl.ID // channel
	Y    chnl.ID // value
	Cont Term
}

func (RecvSpec) cont() {}

func (s RecvSpec) VID() chnl.ID          { return s.X }
func (s RecvSpec) Via() core.Placeholder { return s.X }

type CloseSpec struct {
	A core.Placeholder
}

func (CloseSpec) val() {}

func (s CloseSpec) VID() chnl.ID          { return s.A.(chnl.ID) }
func (s CloseSpec) Via() core.Placeholder { return s.A }

type WaitSpec struct {
	X    chnl.ID
	Cont Term
}

func (WaitSpec) cont() {}

func (s WaitSpec) VID() chnl.ID          { return s.X }
func (s WaitSpec) Via() core.Placeholder { return s.X }

type CTASpec struct {
	Seat sym.ADT
	Key  ak.ADT
}

func (s CTASpec) VID() chnl.ID          { return id.New() }
func (s CTASpec) Via() core.Placeholder { return id.New() }

// aka ExpName
type RecurSpec struct {
	Seat sym.ADT
	C    chnl.ID
	Ctx  []chnl.ID
}

func (s RecurSpec) VID() chnl.ID          { return s.C }
func (s RecurSpec) Via() core.Placeholder { return s.C }

type FwdSpec struct {
	A chnl.ID // from
	C chnl.ID // to
}

func (s FwdSpec) VID() chnl.ID          { return s.C }
func (s FwdSpec) Via() core.Placeholder { return s.C }

type SpawnSpec struct {
	Z      core.Placeholder
	Ctx    map[chnl.Name]chnl.ID
	Cont   Term
	SeatID id.ADT
}

func (s SpawnSpec) VID() chnl.ID          { return s.Z.(chnl.ID) }
func (s SpawnSpec) Via() core.Placeholder { return s.Z }

type Repo[T Root] interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (*T, error)
	SelectByPID(chnl.ID) (*T, error)
	SelectByVID(chnl.ID) (*T, error)
}

// TODO собирать в set, чтобы не было дублей
func CollectChnlIDs(t Term, ids []chnl.ID) []chnl.ID {
	switch term := t.(type) {
	case CloseSpec:
		id, ok := term.A.(chnl.ID)
		if !ok {
			return ids
		}
		return append(ids, id)
	case WaitSpec:
		return CollectChnlIDs(term.Cont, append(ids, term.X))
	case SendSpec:
		return append(ids, term.A, term.B)
	case RecvSpec:
		return CollectChnlIDs(term.Cont, append(ids, term.X, term.Y))
	case LabSpec:
		return append(ids, term.C)
	case CaseSpec:
		for _, cont := range term.Conts {
			ids = CollectChnlIDs(cont, ids)
		}
		return append(ids, term.Z)
	case SpawnSpec:
		ids = CollectChnlIDs(term.Cont, ids)
		id, ok := term.Z.(chnl.ID)
		if !ok {
			return ids
		}
		return append(ids, id)
	default:
		panic(ErrUnexpectedTerm(t))
	}
}

func SubstByID(t Term, varID chnl.ID, valID chnl.ID) Term {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case WaitSpec:
		if varID == term.X {
			term.X = valID
		}
		term.Cont = SubstByID(term.Cont, varID, valID)
		return term
	default:
		panic(ErrUnexpectedTerm(t))
	}
}

func SubstByPH(t Term, ph core.Placeholder, val chnl.ID) Term {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case CloseSpec:
		if ph == term.A {
			term.A = val
		}
		return term
	default:
		panic(ErrUnexpectedTerm(t))
	}
}

func ErrUnexpectedStep(s Root) error {
	return fmt.Errorf("unexpected step type: %T", s)
}

func ErrUnexpectedTerm(t Term) error {
	return fmt.Errorf("unexpected term type: %T", t)
}

func ErrTermMismatch(got, want Term) error {
	return fmt.Errorf("term mismatch: want %T, got %T", want, got)
}

func ErrUnexpectedValue(v Value) error {
	return fmt.Errorf("unexpected value type: %T", v)
}

func ErrUnexpectedCont(c Continuation) error {
	return fmt.Errorf("unexpected continuation type: %T", c)
}

func ErrDoesNotExist(rid ID) error {
	return fmt.Errorf("step doesn't exist: %v", rid)
}
