package state

import (
	"errors"

	"smecalculus/rolevod/lib/id"
)

type ID interface{}

type Ref interface {
	ID() id.ADT[ID]
}

type ref struct {
	id id.ADT[ID]
}

func (r ref) ID() id.ADT[ID] { return r.id }

type WithRef struct{ ref }
type PlusRef struct{ ref }
type TensorRef struct{ ref }
type LolliRef struct{ ref }
type OneRef struct{ ref }
type TpRefRef struct{ ref }
type UpRef struct{ ref }
type DownRef struct{ ref }

// aka Stype
type Root interface {
	getID() id.ADT[ID]
}

type Label string

// aka External Choice
type With struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (r With) getID() id.ADT[ID] { return r.ID }

// aka Internal Choice
type Plus struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (r Plus) getID() id.ADT[ID] { return r.ID }

type Tensor struct {
	ID id.ADT[ID]
	S  Root
	T  Root
}

func (r Tensor) getID() id.ADT[ID] { return r.ID }

type Lolli struct {
	ID id.ADT[ID]
	S  Root
	T  Root
}

func (r Lolli) getID() id.ADT[ID] { return r.ID }

type One struct {
	ID id.ADT[ID]
}

func (r One) getID() id.ADT[ID] { return r.ID }

// TODO тут ссылка на role?
// aka TpName
type TpRef struct {
	ID   id.ADT[ID]
	Name string
}

func (r TpRef) getID() id.ADT[ID] { return r.ID }

type Up struct {
	ID id.ADT[ID]
	A  Root
}

func (r Up) GetID() id.ADT[ID] { return r.ID }

type Down struct {
	ID id.ADT[ID]
	A  Root
}

func (r Down) GetID() id.ADT[ID] { return r.ID }

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT[ID]) (Root, error)
}

var (
	ErrUnexpectedState = errors.New("unexpected state type")
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToCoreIDs func([]string) ([]id.ADT[ID], error)
	ToEdgeIDs func([]id.ADT[ID]) []string
)

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
