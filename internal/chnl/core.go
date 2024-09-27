package chnl

import (
	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type ID = id.ADT

// Symbol
type Sym string

// aka ChanTp
type Spec struct {
	Name Sym
	St   state.Ref
}

// aka Z
type Ref struct {
	ID   ID
	Name Sym
}

// Communication Ep
type Ep struct {
	ID ID
	AK ak.ADT
}

type Root struct {
	ID   ID
	Name Sym
	// Preceding Channel ID
	PreID ID
	// Producer Access Key
	PAK ak.ADT
	// Consumer Access Key
	CAK ak.ADT
	// State
	St state.Ref
}

type Repo interface {
	Insert(Root) error
	InsertCtx([]Root) ([]Root, error)
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Ident
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:Ident
var (
	ConvertRootToRef func(Root) Ref
)
