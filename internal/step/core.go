package step

import (
	"fmt"

	"smecalculus/rolevod/app/seat"
	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
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

type CTASpec struct {
	Key  ak.ADT
	Seat seat.ID
}

func (s CTASpec) Via() core.Placeholder { return id.New() }

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
	Ctx    map[chnl.Name]chnl.ID
	Cont   Term
	SeatID id.ADT
}

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
		ids := CollectChnlIDs(term.Cont, ids)
		id, ok := term.X.(chnl.ID)
		if !ok {
			return ids
		}
		return append(ids, id)
	case SendSpec:
		aID, ok := term.A.(chnl.ID)
		if ok {
			ids = append(ids, aID)
		}
		bID, ok := term.B.(chnl.ID)
		if ok {
			ids = append(ids, bID)
		}
		return ids
	case RecvSpec:
		xID, ok := term.X.(chnl.ID)
		if ok {
			ids = append(ids, xID)
		}
		yID, ok := term.Y.(chnl.ID)
		if ok {
			ids = append(ids, yID)
		}
		return CollectChnlIDs(term.Cont, ids)
	case LabSpec:
		id, ok := term.C.(chnl.ID)
		if !ok {
			return ids
		}
		return append(ids, id)
	case CaseSpec:
		for _, cont := range term.Conts {
			ids = CollectChnlIDs(cont, ids)
		}
		id, ok := term.Z.(chnl.ID)
		if !ok {
			return ids
		}
		return append(ids, id)
	case SpawnSpec:
		ids := CollectChnlIDs(term.Cont, ids)
		id, ok := term.Z.(chnl.ID)
		if !ok {
			return ids
		}
		return append(ids, id)
	case FwdSpec:
		cID, ok := term.C.(chnl.ID)
		if ok {
			ids = append(ids, cID)
		}
		dID, ok := term.D.(chnl.ID)
		if ok {
			ids = append(ids, dID)
		}
		return ids
	default:
		panic(ErrUnexpectedTermType(t))
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
	default:
		panic(ErrUnexpectedTermType(t))
	}
}

func ErrUnexpectedRootType(s Root) error {
	return fmt.Errorf("unexpected root type: %T", s)
}

func ErrUnexpectedTermType(t Term) error {
	return fmt.Errorf("unexpected term type: %T", t)
}

func ErrTermMismatch(got, want Term) error {
	return fmt.Errorf("term mismatch: want %T, got %T", want, got)
}

func ErrUnexpectedValueType(v Value) error {
	return fmt.Errorf("unexpected value type: %T", v)
}

func ErrUnexpectedContType(c Continuation) error {
	return fmt.Errorf("unexpected continuation type: %T", c)
}

func ErrDoesNotExist(rid ID) error {
	return fmt.Errorf("step doesn't exist: %v", rid)
}
