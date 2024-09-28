package chnl

import (
	"database/sql"
	"encoding/json"

	"smecalculus/rolevod/internal/state"
)

type SpecData struct {
	Name string         `json:"name"`
	St   *state.RefData `json:"state"`
}

type RefData struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type rootData struct {
	ID    string         `db:"id"`
	Name  string         `db:"name"`
	PreID string         `db:"pre_id"`
	St    *state.RefData `db:"state"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:String.*
// goverter:extend smecalculus/rolevod/internal/state:Json.*
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

func JsonFromSpec(spec Spec) (sql.NullString, error) {
	dto, err := DataFromSpec(spec)
	if err != nil {
		return sql.NullString{}, err
	}
	jsn, err := json.Marshal(dto)
	if err != nil {
		return sql.NullString{}, err
	}
	return sql.NullString{String: string(jsn), Valid: true}, nil
}

func JsonToSpec(jsn sql.NullString) (Spec, error) {
	if !jsn.Valid {
		return Spec{}, nil
	}
	var dto SpecData
	err := json.Unmarshal([]byte(jsn.String), &dto)
	if err != nil {
		return Spec{}, err
	}
	return DataToSpec(dto)
}

func JsonFromSpecs(specs []Spec) (sql.NullString, error) {
	dtos, err := DataFromSpecs(specs)
	if err != nil {
		return sql.NullString{}, err
	}
	jsn, err := json.Marshal(dtos)
	if err != nil {
		return sql.NullString{}, err
	}
	return sql.NullString{String: string(jsn), Valid: true}, nil
}

func JsonToSpecs(jsn sql.NullString) ([]Spec, error) {
	if !jsn.Valid {
		return []Spec{}, nil
	}
	var dtos []SpecData
	err := json.Unmarshal([]byte(jsn.String), &dtos)
	if err != nil {
		return nil, err
	}
	return DataToSpecs(dtos)
}
