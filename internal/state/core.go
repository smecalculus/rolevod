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

type WithRef id.ADT[ID]

func (r WithRef) RID() id.ADT[ID] { return id.ADT[ID](r) }

type PlusRef id.ADT[ID]

func (r PlusRef) RID() id.ADT[ID] { return id.ADT[ID](r) }

type TensorRef id.ADT[ID]

func (r TensorRef) RID() id.ADT[ID] { return id.ADT[ID](r) }

type LolliRef id.ADT[ID]

func (r LolliRef) RID() id.ADT[ID] { return id.ADT[ID](r) }

type OneRef id.ADT[ID]

func (r OneRef) RID() id.ADT[ID] { return id.ADT[ID](r) }

type RecurRef id.ADT[ID]

func (r RecurRef) RID() id.ADT[ID] { return id.ADT[ID](r) }

type UpRef id.ADT[ID]

func (r UpRef) RID() id.ADT[ID] { return id.ADT[ID](r) }

type DownRef id.ADT[ID]

func (r DownRef) RID() id.ADT[ID] { return id.ADT[ID](r) }

// aka Stype
type Root interface {
	Spec
	Ref
}

// aka External Choice
type WithRoot struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (WithRoot) spec() {}

func (r WithRoot) RID() id.ADT[ID] { return r.ID }

// aka Internal Choice
type PlusRoot struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (PlusRoot) spec() {}

func (r PlusRoot) RID() id.ADT[ID] { return r.ID }

type TensorRoot struct {
	ID id.ADT[ID]
	A  Root // value
	C  Root // cont
}

func (TensorRoot) spec() {}

func (r TensorRoot) RID() id.ADT[ID] { return r.ID }

type LolliRoot struct {
	ID id.ADT[ID]
	X  Root // value
	Z  Root // cont
}

func (LolliRoot) spec() {}

func (r LolliRoot) RID() id.ADT[ID] { return r.ID }

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
			A:  ConvertSpecToRoot(spec.A),
			C:  ConvertSpecToRoot(spec.C),
		}
	case LolliSpec:
		return LolliRoot{
			ID: newID,
			X:  ConvertSpecToRoot(spec.X),
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
	switch root := r.(type) {
	case OneRoot:
		return OneRef(root.ID)
	case RecurRoot:
		return RecurRef(root.ID)
	case TensorRoot:
		return TensorRef(root.ID)
	case LolliRoot:
		return LolliRef(root.ID)
	case WithRoot:
		return WithRef(root.ID)
	case PlusRoot:
		return PlusRef(root.ID)
	default:
		panic(ErrUnexpectedRoot(r))
	}
}

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
