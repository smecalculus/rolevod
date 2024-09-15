package state

import (
	"errors"
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

func (r WithSpec) spec() {}

// aka Internal Choice
type PlusSpec struct {
	Choices map[Label]Spec
}

func (r PlusSpec) spec() {}

type TensorSpec struct {
	S Spec
	T Spec
}

func (r TensorSpec) spec() {}

type LolliSpec struct {
	S Spec
	T Spec
}

func (r LolliSpec) spec() {}

type OneSpec struct{}

func (r OneSpec) spec() {}

// aka TpName
type TpRefSpec struct {
	ID   id.ADT[ID]
	Name string
}

func (r TpRefSpec) spec() {}

type UpSpec struct {
	A Spec
}

func (r UpSpec) spec() {}

type DownSpec struct {
	A Spec
}

func (r DownSpec) spec() {}

type ID interface{}

type Ref interface {
	RootID() id.ADT[ID]
}

// type ref struct {
// 	id id.ADT[ID]
// }

// func (r ref) rootID() id.ADT[ID] { return r.id }

// type WithRef struct{ ref }
// type PlusRef struct{ ref }
// type TensorRef struct{ ref }
// type LolliRef struct{ ref }
// type OneRef struct{ ref }
// type TpRefRef struct{ ref }
// type UpRef struct{ ref }
// type DownRef struct{ ref }

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

type TpRefRef id.ADT[ID]

func (r TpRefRef) RootID() id.ADT[ID] { return id.ADT[ID](r) }

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

// TODO тут ссылка на role?
// aka TpName
type TpRefRoot struct {
	ID      id.ADT[ID]
	Name    string
	StateID id.ADT[ID]
}

func (r TpRefRoot) rootID() id.ADT[ID] { return r.ID }

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

var (
	ErrUnexpectedState = errors.New("unexpected state type")
)

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
	switch spec := s.(type) {
	case OneSpec:
		// TODO генерировать zero id или не генерировать id вообще
		return OneRoot{ID: id.New[ID]()}
	case TpRefSpec:
		return TpRefRoot{ID: spec.ID, Name: spec.Name}
	case TensorSpec:
		return TensorRoot{
			ID: id.New[ID](),
			S:  ConvertSpecToRoot(spec.S),
			T:  ConvertSpecToRoot(spec.T),
		}
	case LolliSpec:
		return LolliRoot{
			ID: id.New[ID](),
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
		// return OneRef{ref{root.ID}}
		return OneRef(root.ID)
	case TpRefRoot:
		// return TpRefRef{ref{root.ID}}
		return TpRefRef(root.ID)
	case TensorRoot:
		// return TensorRef{ref{root.ID}}
		return TensorRef(root.ID)
	case LolliRoot:
		// return LolliRef{ref{root.ID}}
		return LolliRef(root.ID)
	case WithRoot:
		// return WithRef{ref{root.ID}}
		return WithRef(root.ID)
	case PlusRoot:
		// return PlusRef{ref{root.ID}}
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
