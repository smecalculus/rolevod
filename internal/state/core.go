package state

import (
	"fmt"

	"smecalculus/rolevod/lib/id"
)

type Label string

type Spec interface {
	spec()
}

// aka External Choice
type WithSpec struct {
	Choices map[Label]Spec
}

func (WithSpec) spec() {}

// aka Internal Choice
type PlusSpec struct {
	Choices map[Label]Spec
}

func (PlusSpec) spec() {}

type TensorSpec struct {
	A Spec
	C Spec
}

func (TensorSpec) spec() {}

type LolliSpec struct {
	X Spec
	Z Spec
}

func (LolliSpec) spec() {}

type OneSpec struct{}

func (OneSpec) spec() {}

// aka TpName
type RecurSpec struct {
	Name string
	ToID id.ADT[ID]
}

func (RecurSpec) spec() {}

type UpSpec struct {
	A Spec
}

func (UpSpec) spec() {}

type DownSpec struct {
	A Spec
}

func (DownSpec) spec() {}

type ID interface{}

type Ref interface {
	RID() id.ADT[ID]
}

type WithRef struct {
	ID id.ADT[ID]
}

func (r WithRef) RID() id.ADT[ID] { return r.ID }

type PlusRef struct {
	ID id.ADT[ID]
}

func (r PlusRef) RID() id.ADT[ID] { return r.ID }

type TensorRef struct {
	ID id.ADT[ID]
}

func (r TensorRef) RID() id.ADT[ID] { return r.ID }

type LolliRef struct {
	ID id.ADT[ID]
}

func (r LolliRef) RID() id.ADT[ID] { return r.ID }

type OneRef struct {
	ID id.ADT[ID]
}

func (r OneRef) RID() id.ADT[ID] { return r.ID }

type RecurRef struct {
	ID id.ADT[ID]
}

func (r RecurRef) RID() id.ADT[ID] { return r.ID }

type UpRef struct {
	ID id.ADT[ID]
}

func (r UpRef) RID() id.ADT[ID] { return r.ID }

type DownRef struct {
	ID id.ADT[ID]
}

func (r DownRef) RID() id.ADT[ID] { return r.ID }

// aka Stype
type Root interface {
	Ref
}

type Product interface {
	Next() Ref
}

type Sum interface {
	Next(Label) Ref
}

// aka Internal Choice
type PlusRoot struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (r PlusRoot) RID() id.ADT[ID] { return r.ID }

func (r PlusRoot) Next(l Label) Ref { return r.Choices[l] }

// aka External Choice
type WithRoot struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (r WithRoot) RID() id.ADT[ID] { return r.ID }

func (r WithRoot) Next(l Label) Ref { return r.Choices[l] }

type TensorRoot struct {
	ID id.ADT[ID]
	B  Ref  // value
	C  Root // cont
}

func (r TensorRoot) RID() id.ADT[ID] { return r.ID }

func (r TensorRoot) Next() Ref { return r.C }

type LolliRoot struct {
	ID id.ADT[ID]
	Y  Ref  // value
	Z  Root // cont
}

func (r LolliRoot) RID() id.ADT[ID] { return r.ID }

func (r LolliRoot) Next() Ref { return r.Z }

type OneRoot struct {
	ID id.ADT[ID]
}

func (OneRoot) spec() {}

func (r OneRoot) RID() id.ADT[ID] { return r.ID }

// aka TpName
type RecurRoot struct {
	ID   id.ADT[ID]
	Name string
	ToID id.ADT[ID]
}

func (RecurRoot) spec() {}

func (r RecurRoot) RID() id.ADT[ID] { return r.ID }

type UpRoot struct {
	ID id.ADT[ID]
	A  Root
}

func (UpRoot) spec() {}

func (r UpRoot) RID() id.ADT[ID] { return r.ID }

type DownRoot struct {
	ID id.ADT[ID]
	A  Root
}

func (DownRoot) spec() {}

func (r DownRoot) RID() id.ADT[ID] { return r.ID }

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT[ID]) (Root, error)
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

func ConvertSpecToRoot(s Spec) Root {
	if s == nil {
		return nil
	}
	newID := id.New[ID]()
	switch spec := s.(type) {
	case OneSpec:
		// TODO генерировать zero id или не генерировать id вообще
		return OneRoot{ID: newID}
	case RecurSpec:
		return RecurRoot{ID: newID, Name: spec.Name, ToID: spec.ToID}
	case TensorSpec:
		return TensorRoot{
			ID: newID,
			B:  ConvertSpecToRoot(spec.A),
			C:  ConvertSpecToRoot(spec.C),
		}
	case LolliSpec:
		return LolliRoot{
			ID: newID,
			Y:  ConvertSpecToRoot(spec.X),
			Z:  ConvertSpecToRoot(spec.Z),
		}
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

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
