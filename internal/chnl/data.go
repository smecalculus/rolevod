package chnl

import (
	"smecalculus/rolevod/internal/state"
)

type SpecData struct {
	Name string         `json:"name"`
	StID string         `json:"st_id"`
	St   *state.RefData `json:"state"`
}

type RefData struct {
	ID   string `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

type rootData struct {
	ID    string         `db:"id"`
	Name  string         `db:"name"`
	PreID string         `db:"pre_id"` // TODO null string
	StID  string         `db:"st_id"`  // TODO null string
	St    *state.RefData `db:"state"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:String.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
var (
	DataToSpec    func(SpecData) (Spec, error)
	DataFromSpec  func(Spec) (SpecData, error)
	DataToSpecs   func([]SpecData) ([]Spec, error)
	DataFromSpecs func([]Spec) ([]SpecData, error)
	DataToRef     func(RefData) (Ref, error)
	DataFromRef   func(Ref) RefData
	DataToRefs    func([]RefData) ([]Ref, error)
	DataFromRefs  func([]Ref) []RefData
	DataToRoot    func(rootData) (Root, error)
	DataFromRoot  func(Root) (rootData, error)
	DataToRoots   func([]rootData) ([]Root, error)
	DataFromRoots func([]Root) ([]rootData, error)
)
