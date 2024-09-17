package chnl

import (
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type ID interface{}

// aka ChanTp
type Spec struct {
	Name string
	St   state.Ref
}

// aka Z
type Ref struct {
	ID   id.ADT[ID]
	Name string
}

type Var string

// aka ChanTp with ID
type Root struct {
	ID    id.ADT[ID]
	Name  string
	PreID id.ADT[ID]
	St    state.Ref
}

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT[ID]) (Root, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ConvertRootToRef func(Root) Ref
)

func toSame(id id.ADT[ID]) id.ADT[ID] {
	return id
}

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
