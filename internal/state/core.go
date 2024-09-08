package state

import (
	"errors"

	"smecalculus/rolevod/lib/id"
)

type ID interface{}

type Ref interface {
	Get() id.ADT[ID]
}

type Label string

// aka Stype
type Root interface {
	Ref
}

// aka External Choice
type With struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (r With) Get() id.ADT[ID] { return r.ID }

// aka Internal Choice
type Plus struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (r Plus) Get() id.ADT[ID] { return r.ID }

type Tensor struct {
	ID id.ADT[ID]
	S  Root
	T  Root
}

func (r Tensor) Get() id.ADT[ID] { return r.ID }

type Lolli struct {
	ID id.ADT[ID]
	S  Root
	T  Root
}

func (r Lolli) Get() id.ADT[ID] { return r.ID }

type One struct {
	ID id.ADT[ID]
}

func (r One) Get() id.ADT[ID] { return r.ID }

// aka TpName
type TpRef struct {
	ID id.ADT[ID]
	// TODO: выпилить
	Name string
}

func (r TpRef) Get() id.ADT[ID] { return r.ID }

type Up struct {
	ID id.ADT[ID]
	A  Root
}

func (r Up) Get() id.ADT[ID] { return r.ID }

type Down struct {
	ID id.ADT[ID]
	A  Root
}

func (r Down) Get() id.ADT[ID] { return r.ID }

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectById(id.ADT[ID]) (Root, error)
	SelectNext(id.ADT[ID]) (Ref, error)
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
