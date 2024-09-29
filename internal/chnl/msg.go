package chnl

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type SpecMsg struct {
	Name string       `json:"name"`
	St   state.RefMsg `json:"state"`
}

func (mto SpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Name, core.NameRequired...),
		validation.Field(&mto.St, validation.Required),
	)
}

type RefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

func (mto RefMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.ID, id.Required...),
		validation.Field(&mto.Name, core.NameRequired...),
	)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:String.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToSpec   func(SpecMsg) (Spec, error)
	MsgFromSpec func(Spec) SpecMsg
	MsgToRef    func(RefMsg) (Ref, error)
	MsgFromRef  func(Ref) RefMsg
)

func MsgFromRefMap(refs map[Name]ID) []RefMsg {
	var mtos []RefMsg
	for name, rid := range refs {
		mtos = append(mtos, RefMsg{rid.String(), name})
	}
	return mtos
}

func MsgToRefMap(mtos []RefMsg) (map[Name]ID, error) {
	refs := make(map[Name]ID, len(mtos))
	for _, mto := range mtos {
		mtoID, err := id.StringToID(mto.ID)
		if err != nil {
			return nil, err
		}
		refs[mto.Name] = mtoID
	}
	return refs, nil
}
