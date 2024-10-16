package chnl

import (
	"database/sql"

	"smecalculus/rolevod/lib/id"
)

type SpecData struct {
	Name string `json:"name"`
	StID string `json:"st_id"`
}

type RefData struct {
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
	DataToRef     func(RefData) (Ref, error)
	DataFromRef   func(Ref) RefData
	DataToRefs    func([]RefData) ([]Ref, error)
	DataFromRefs  func([]Ref) []RefData
	DataToRoot    func(rootData) (Root, error)
	DataFromRoot  func(Root) (rootData, error)
	DataToRoots   func([]rootData) ([]Root, error)
	DataFromRoots func([]Root) ([]rootData, error)
)

func DataFromRefMap(refs map[Name]ID) []RefData {
	var dtos []RefData
	for name, rid := range refs {
		dtos = append(dtos, RefData{rid.String(), name})
	}
	return dtos
}

func DataToRefMap(dtos []RefData) (map[Name]ID, error) {
	refs := make(map[Name]ID, len(dtos))
	for _, mto := range dtos {
		dtoID, err := id.StringToID(mto.ID)
		if err != nil {
			return nil, err
		}
		refs[mto.Name] = dtoID
	}
	return refs, nil
}
