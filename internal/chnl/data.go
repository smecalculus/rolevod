package chnl

import (
	"database/sql"

	"smecalculus/rolevod/lib/rev"
)

type SpecData struct {
	Key  string `json:"chnl_key"`
	Link string `json:"role_fqn"`
}

type refData struct {
	ID    string `db:"chnl_id" json:"id,omitempty"`
	Title string `db:"title" json:"title,omitempty"`
}

type rootData struct {
	ID      string         `db:"chnl_id"`
	Title   string         `db:"title"`
	PoolID  sql.NullString `db:"pool_id"`
	StateID sql.NullString `db:"state_id"`
	Revs    []rev.ADT      `db:"revs"`
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
