package state

import (
	"errors"

	"smecalculus/rolevod/lib/id"
)

var (
	ErrUnexpectedSt = errors.New("unexpected session type")
)

type ID interface{}

type Ref struct {
	ID id.ADT[ID]
	// TODO: доменное представление kind
}

type Label string

// aka Stype
type Root interface {
	get() id.ADT[ID]
	set(id.ADT[ID])
}

// aka External Choice
type With struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (s *With) get() id.ADT[ID]   { return s.ID }
func (s *With) set(id id.ADT[ID]) { s.ID = id }

// aka Internal Choice
type Plus struct {
	ID      id.ADT[ID]
	Choices map[Label]Root
}

func (s *Plus) get() id.ADT[ID]   { return s.ID }
func (s *Plus) set(id id.ADT[ID]) { s.ID = id }

type Tensor struct {
	ID id.ADT[ID]
	S  Root
	T  Root
}

func (s *Tensor) get() id.ADT[ID]   { return s.ID }
func (s *Tensor) set(id id.ADT[ID]) { s.ID = id }

type Lolli struct {
	ID id.ADT[ID]
	S  Root
	T  Root
}

func (s *Lolli) get() id.ADT[ID]   { return s.ID }
func (s *Lolli) set(id id.ADT[ID]) { s.ID = id }

type One struct {
	ID id.ADT[ID]
}

func (s *One) get() id.ADT[ID]   { return s.ID }
func (s *One) set(id id.ADT[ID]) { s.ID = id }

// TODO: выпилить Name
// aka TpName
type TpRef struct {
	ID   id.ADT[ID]
	Name string
}

func (s *TpRef) get() id.ADT[ID]   { return s.ID }
func (s *TpRef) set(id id.ADT[ID]) { s.ID = id }

type Up struct {
	ID id.ADT[ID]
	A  Root
}

func (s *Up) get() id.ADT[ID]   { return s.ID }
func (s *Up) set(id id.ADT[ID]) { s.ID = id }

type Down struct {
	ID id.ADT[ID]
	A  Root
}

func (s *Down) get() id.ADT[ID]   { return s.ID }
func (s *Down) set(id id.ADT[ID]) { s.ID = id }

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectById(id.ADT[ID]) (Root, error)
}

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
