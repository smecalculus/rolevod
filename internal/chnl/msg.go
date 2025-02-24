package chnl

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type SpecMsg struct {
	Key  string `json:"name"`
	Link string `json:"role_fqn"`
}

func (dto SpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Key, core.NameRequired...),
		validation.Field(&dto.Link, sym.Required...),
	)
}

type RefMsg struct {
	ID  string `json:"id" param:"id"`
	Key string `json:"name"`
}

func (dto RefMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
		validation.Field(&dto.Key, core.NameRequired...),
	)
}

type RootMsg struct {
	ID      string  `json:"id" param:"id"`
	Key     string  `json:"name"`
	PreID   *string `json:"pre_id"`
	StateID *string `json:"state_id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/ak:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToSpec   func(SpecMsg) (Spec, error)
	MsgFromSpec func(Spec) SpecMsg
	MsgToRef    func(RefMsg) (Ref, error)
	MsgFromRef  func(Ref) RefMsg
	MsgToRoot   func(RootMsg) (Root, error)
	MsgFromRoot func(Root) RootMsg
)

func MsgFromRefMap(refs map[Key]ID) []RefMsg {
	var mtos []RefMsg
	for key, rid := range refs {
		mtos = append(mtos, RefMsg{rid.String(), key})
	}
	return mtos
}

func MsgToRefMap(mtos []RefMsg) (map[Key]ID, error) {
	refs := make(map[Key]ID, len(mtos))
	for _, dto := range mtos {
		mtoID, err := id.ConvertFromString(dto.ID)
		if err != nil {
			return nil, err
		}
		refs[dto.Key] = mtoID
	}
	return refs, nil
}
