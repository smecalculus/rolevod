package chnl

import (
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

type Root struct {
	ID   ID
	Name Sym
	// Preceding Channel ID
	PreID ID
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
