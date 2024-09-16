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
	S Spec
	T Spec
}

func (TensorSpec) spec() {}

type LolliSpec struct {
	S Spec
	T Spec
}

func (LolliSpec) spec() {}

type OneSpec struct{}

func (OneSpec) spec() {}

// aka TpName
type RecSpec struct {
	Name string
	ToID id.ADT[ID]
}

func (RecSpec) spec() {}

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
	RootID() id.ADT[ID]
}

type WithRef id.ADT[ID]

func (r WithRef) RootID() id.ADT[ID] { return id.ADT[ID](r) }

type PlusRef id.ADT[ID]

func (r PlusRef) RootID() id.ADT[ID] { return id.ADT[ID](r) }

type TensorRef id.ADT[ID]

func (r TensorRef) RootID() id.ADT[ID] { return id.ADT[ID](r) }

type LolliRef id.ADT[ID]

func (r LolliRef) RootID() id.ADT[ID] { return id.ADT[ID](r) }

type OneRef id.ADT[ID]

func (r OneRef) RootID() id.ADT[ID] { return id.ADT[ID](r) }

type RecRef id.ADT[ID]

func (r RecRef) RootID() id.ADT[ID] { return id.ADT[ID](r) }

type UpRef id.ADT[ID]

func (r UpRef) rootID() id.ADT[ID] { return id.ADT[ID](r) }

type DownRef id.ADT[ID]

func (r DownRef) rootID() id.ADT[ID] { return id.ADT[ID](r) }

// aka Stype
type Root interface {
	rootID() id.ADT[ID]
}

// aka External Choice
type WithRoot struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (r WithRoot) rootID() id.ADT[ID] { return r.ID }

// aka Internal Choice
type PlusRoot struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (r PlusRoot) rootID() id.ADT[ID] { return r.ID }

type TensorRoot struct {
	ID id.ADT[ID]
	S  Root
	T  Root
}

func (r TensorRoot) rootID() id.ADT[ID] { return r.ID }

type LolliRoot struct {
	ID id.ADT[ID]
	S  Root
	T  Root
}

func (r LolliRoot) rootID() id.ADT[ID] { return r.ID }

type OneRoot struct {
	ID id.ADT[ID]
}

func (r OneRoot) rootID() id.ADT[ID] { return r.ID }

// aka TpName
type RecRoot struct {
	ID   id.ADT[ID]
	Name string
	ToID id.ADT[ID]
}

func (r RecRoot) rootID() id.ADT[ID] { return r.ID }

type UpRoot struct {
	ID id.ADT[ID]
	A  Root
}

func (r UpRoot) rootID() id.ADT[ID] { return r.ID }

type DownRoot struct {
	ID id.ADT[ID]
	A  Root
}

func (r DownRoot) rootID() id.ADT[ID] { return r.ID }

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
	case RecSpec:
		return RecRoot{ID: newID, Name: spec.Name, ToID: spec.ToID}
	case TensorSpec:
		return TensorRoot{
			ID: newID,
			S:  ConvertSpecToRoot(spec.S),
			T:  ConvertSpecToRoot(spec.T),
		}
	case LolliSpec:
		return LolliRoot{
			ID: newID,
			S:  ConvertSpecToRoot(spec.S),
			T:  ConvertSpecToRoot(spec.T),
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
	case RecRoot:
		return RecRef(root.ID)
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
