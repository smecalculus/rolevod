package chnl

import (
	"encoding/json"

	"smecalculus/rolevod/internal/state"
)

type refData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type rootData struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	State string `db:"state"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend state.*
var (
	DataToRef    func(refData) (Ref, error)
	DataFromRef  func(Ref) refData
	DataToRefs   func([]refData) ([]Ref, error)
	DataFromRefs func([]Ref) []refData
	DataToRoot   func(rootData) (Root, error)
	DataFromRoot func(Root) rootData
)

func stateDataFromRef(ref state.Ref) (string, error) {
	dto := state.DataFromRef(ref)
	s, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func stateDataToRef(data string) (state.Ref, error) {
	var dto state.RefData
	err := json.Unmarshal([]byte(data), &dto)
	if err != nil {
		return nil, err
	}
	return state.DataToRef(dto)
}
