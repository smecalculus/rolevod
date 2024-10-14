package state

import (
	"fmt"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type ID = id.ADT

type Spec interface {
	spec()
}

type OneSpec struct{}

func (OneSpec) spec() {}

// aka TpName
type LinkSpec struct {
	FQN sym.ADT
}

func (LinkSpec) spec() {}

type TensorSpec struct {
	B Spec
	C Spec
}

func (TensorSpec) spec() {}

type LolliSpec struct {
	Y Spec
	Z Spec
}

func (LolliSpec) spec() {}

// aka Internal Choice
type PlusSpec struct {
	Choices map[core.Label]Spec
}

func (PlusSpec) spec() {}

// aka External Choice
type WithSpec struct {
	Choices map[core.Label]Spec
}

func (WithSpec) spec() {}

type UpSpec struct {
	A Spec
}

func (UpSpec) spec() {}

type DownSpec struct {
	A Spec
}

func (DownSpec) spec() {}

type Ref interface {
	RID() ID
}

type OneRef struct {
	ID ID
}

func (r OneRef) RID() ID { return r.ID }

type LinkRef struct {
	ID ID
}

func (r LinkRef) RID() ID { return r.ID }

type PlusRef struct {
	ID ID
}

func (r PlusRef) RID() ID { return r.ID }

type WithRef struct {
	ID ID
}

func (r WithRef) RID() ID { return r.ID }

type TensorRef struct {
	ID ID
}

func (r TensorRef) RID() ID { return r.ID }

type LolliRef struct {
	ID ID
}

func (r LolliRef) RID() ID { return r.ID }

type UpRef struct {
	ID ID
}

func (r UpRef) RID() ID { return r.ID }

type DownRef struct {
	ID ID
}

func (r DownRef) RID() ID { return r.ID }

// aka Stype
type Root interface {
	Ref
	Pol() Polarity
}

type Prod interface {
	Next() Ref
}

type Sum interface {
	Next(core.Label) Ref
}

// aka TpName
type LinkRoot struct {
	ID  ID
	FQN sym.ADT
}

func (LinkRoot) spec() {}

func (r LinkRoot) RID() ID { return r.ID }

func (r LinkRoot) Pol() Polarity { return Zero }

type OneRoot struct {
	ID ID
}

func (OneRoot) spec() {}

func (r OneRoot) RID() ID { return r.ID }

func (r OneRoot) Pol() Polarity { return Pos }

// aka Internal Choice
type PlusRoot struct {
	ID      ID
	Choices map[core.Label]Root
}

func (r PlusRoot) RID() ID { return r.ID }

func (r PlusRoot) Next(l core.Label) Ref { return r.Choices[l] }

func (r PlusRoot) Pol() Polarity { return Pos }

// aka External Choice
type WithRoot struct {
	ID      ID
	Choices map[core.Label]Root
}

func (r WithRoot) RID() ID { return r.ID }

func (r WithRoot) Next(l core.Label) Ref { return r.Choices[l] }

func (r WithRoot) Pol() Polarity { return Neg }

type TensorRoot struct {
	ID ID
	B  Root // value
	C  Root // cont
}

func (r TensorRoot) RID() ID { return r.ID }

func (r TensorRoot) Next() Ref { return r.C }

func (r TensorRoot) Pol() Polarity { return Pos }

type LolliRoot struct {
	ID ID
	Y  Root // value
	Z  Root // cont
}

func (r LolliRoot) RID() ID { return r.ID }

func (r LolliRoot) Next() Ref { return r.Z }

func (r LolliRoot) Pol() Polarity { return Neg }

type UpRoot struct {
	ID ID
	A  Root
}

func (UpRoot) spec() {}

func (r UpRoot) RID() ID { return r.ID }

func (r UpRoot) Pol() Polarity { return Zero }

type DownRoot struct {
	ID ID
	A  Root
}

func (DownRoot) spec() {}

func (r DownRoot) RID() ID { return r.ID }

func (r DownRoot) Pol() Polarity { return Zero }

type Polarity int

const (
	Pos  = Polarity(+1)
	Zero = Polarity(0)
	Neg  = Polarity(-1)
)

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectEnv([]ID) (map[ID]Root, error)
	SelectMany([]ID) ([]Root, error)
}

func ConvertSpecToRoot(s Spec) Root {
	if s == nil {
		return nil
	}
	newID := id.New()
	switch spec := s.(type) {
	case OneSpec:
		// TODO генерировать zero id или не генерировать id вообще
		return OneRoot{ID: newID}
	case LinkSpec:
		return LinkRoot{ID: newID, FQN: spec.FQN}
	case TensorSpec:
		return TensorRoot{
			ID: newID,
			B:  ConvertSpecToRoot(spec.B),
			C:  ConvertSpecToRoot(spec.C),
		}
	case LolliSpec:
		return LolliRoot{
			ID: newID,
			Y:  ConvertSpecToRoot(spec.Y),
			Z:  ConvertSpecToRoot(spec.Z),
		}
	case WithSpec:
		choices := make(map[core.Label]Root, len(spec.Choices))
		for lab, st := range spec.Choices {
			choices[lab] = ConvertSpecToRoot(st)
		}
		return WithRoot{ID: newID, Choices: choices}
	case PlusSpec:
		choices := make(map[core.Label]Root, len(spec.Choices))
		for lab, st := range spec.Choices {
			choices[lab] = ConvertSpecToRoot(st)
		}
		return PlusRoot{ID: newID, Choices: choices}
	default:
		panic(ErrUnexpectedSpecType(spec))
	}
}

func ErrUnexpectedSpecType(got Spec) error {
	return fmt.Errorf("spec type unexpected: %T", got)
}

func ErrUnexpectedRefType(got Ref) error {
	return fmt.Errorf("ref type unexpected: %T", got)
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
