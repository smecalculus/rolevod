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
	ID   id.ADT
	Name string
	PAK  ak.ADT
	CAK  ak.ADT
}

// Communication Endpoint
type EP struct {
	ID id.ADT
	AK ak.ADT
}

// Symbol
type Sym string

type Root struct {
	ID   id.ADT
	Name string
	// Preceding Channel ID
	PreID id.ADT
	// Producer Access Key
	PAK ak.ADT
	// Consumer Access Key
	CAK ak.ADT
	St  state.Ref
}

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT) (Root, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Ident
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:Ident
var (
	ConvertRootToRef func(Root) Ref
)
