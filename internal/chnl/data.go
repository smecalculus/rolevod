package chnl

import (
	"database/sql"
)

type SpecData struct {
	Name string `json:"name"`
	Role string `json:"role_id"`
}

type refData struct {
	ID   string `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

type rootData struct {
	ID      string         `db:"id"`
	Name    string         `db:"name"`
	PreID   sql.NullString `db:"pre_id"`
	StateID sql.NullString `db:"state_id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/ak:Convert.*
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
