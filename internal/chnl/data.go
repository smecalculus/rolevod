package chnl

import (
	"database/sql"
)

type SpecData struct {
	Name string `json:"name"`
	StID string `json:"st_id"`
}

type refData struct {
	ID   string `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

type rootData struct {
	ID    string         `db:"id"`
	Name  string         `db:"name"`
	PreID sql.NullString `db:"pre_id"`
	StID  sql.NullString `db:"st_id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:String.*
// goverter:extend smecalculus/rolevod/lib/data:NullString.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
var (
	DataToSpec    func(SpecData) (Spec, error)
	DataFromSpec  func(Spec) (SpecData, error)
	DataToSpecs   func([]SpecData) ([]Spec, error)
	DataFromSpecs func([]Spec) ([]SpecData, error)
	DataToRef     func(refData) (Ref, error)
	DataFromRef   func(Ref) refData
	DataToRefs    func([]refData) ([]Ref, error)
	DataFromRefs  func([]Ref) []refData
	DataToRoot    func(rootData) (Root, error)
	DataFromRoot  func(Root) (rootData, error)
	DataToRoots   func([]rootData) ([]Root, error)
	DataFromRoots func([]Root) ([]rootData, error)
)
