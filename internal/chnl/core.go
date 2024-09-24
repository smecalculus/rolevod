package chnl

import (
	"smecalculus/rolevod/lib/ak"
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
	PAK  ak.ADT
	CAK  ak.ADT
}

type Var string

type Root struct {
	ID   id.ADT[ID]
	Name string
	// Preceding Channel ID
	PreID id.ADT[ID]
	// Producer Access Key
	PAK ak.ADT
	// Consumer Access Key
	CAK ak.ADT
	St  state.Ref
}

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT[ID]) (Root, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/lib/ak:Ident.*
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
