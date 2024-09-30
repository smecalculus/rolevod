package state

import (
	"fmt"

	"smecalculus/rolevod/lib/id"
)

type ID = id.ADT

type Label string

type Spec interface {
	spec()
}

type OneSpec struct{}

func (OneSpec) spec() {}

// Mention aka TpName
type MenSpec struct {
	StID ID
	Name string
}

func (MenSpec) spec() {}

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
	Choices map[Label]Spec
}

func (PlusSpec) spec() {}

// aka External Choice
type WithSpec struct {
	Choices map[Label]Spec
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

type MenRef struct {
	ID ID
}

func (r MenRef) RID() ID { return r.ID }

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
}

type Prod interface {
	Next() Ref
}

type Sum interface {
	Next(Label) Ref
}

// aka TpName
type MenRoot struct {
	ID   ID
	StID ID
	Name string
}

func (MenRoot) spec() {}

func (r MenRoot) RID() ID { return r.ID }

type OneRoot struct {
	ID ID
}

func (OneRoot) spec() {}

func (r OneRoot) RID() ID { return r.ID }

// aka Internal Choice
type PlusRoot struct {
	ID      ID
	Choices map[Label]Root
}

func (r PlusRoot) RID() ID { return r.ID }

func (r PlusRoot) Next(l Label) Ref { return r.Choices[l] }

// aka External Choice
type WithRoot struct {
	ID      ID
	Choices map[Label]Root
}

func (r WithRoot) RID() ID { return r.ID }

func (r WithRoot) Next(l Label) Ref { return r.Choices[l] }

type TensorRoot struct {
	ID ID
	B  Root // value
	C  Root // cont
}

func (r TensorRoot) RID() ID { return r.ID }

func (r TensorRoot) Next() Ref { return r.C }

type LolliRoot struct {
	ID ID
	Y  Root // value
	Z  Root // cont
}

func (r LolliRoot) RID() ID { return r.ID }

func (r LolliRoot) Next() Ref { return r.Z }

type UpRoot struct {
	ID ID
	A  Root
}

func (UpRoot) spec() {}

func (r UpRoot) RID() ID { return r.ID }

type DownRoot struct {
	ID ID
	A  Root
}

func (DownRoot) spec() {}

func (r DownRoot) RID() ID { return r.ID }

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
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
	case MenSpec:
		return MenRoot{ID: newID, Name: spec.Name, StID: spec.StID}
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
		choices := make(map[Label]Root, len(spec.Choices))
		for lab, st := range spec.Choices {
			choices[lab] = ConvertSpecToRoot(st)
		}
		return WithRoot{ID: newID, Choices: choices}
	case PlusSpec:
		choices := make(map[Label]Root, len(spec.Choices))
		for lab, st := range spec.Choices {
			choices[lab] = ConvertSpecToRoot(st)
		}
		return PlusRoot{ID: newID, Choices: choices}
	default:
		panic(ErrUnexpectedSpec(spec))
	}
}

func ConvertRefToRef(r Ref) Ref {
	return r
}

func ConvertRootToRef(r Root) Ref {
	if r == nil {
		return nil
	}
	return r.(Ref)
}

func ErrUnexpectedSpec(v Spec) error {
	return fmt.Errorf("unexpected spec %#v", v)
}

func ErrUnexpectedRef(v Ref) error {
	return fmt.Errorf("unexpected ref %#v", v)
}

func ErrUnexpectedRoot(v Root) error {
	return fmt.Errorf("unexpected root %#v", v)
}
