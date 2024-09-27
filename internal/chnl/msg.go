package chnl

import (
	valid "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/internal/state"
)

type SpecMsg struct {
	Name string        `json:"name"`
	St   *state.RefMsg `json:"state"`
}

func (mto *SpecMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.Name, valid.Required, valid.Max(64)),
		valid.Field(&mto.St, valid.Required),
	)
}

type RefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

func (mto *RefMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.ID, valid.Required, valid.Length(20, 20)),
		valid.Field(&mto.Name, valid.Required, valid.Length(1, 64)),
	)
}

type EpMsg struct {
	ID string `json:"id"`
	AK string `json:"ak"`
}

type RootMsg struct {
	ID   string        `json:"id"`
	Name string        `json:"name"`
	St   *state.RefMsg `json:"state"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:String.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromSpec  func(Spec) SpecMsg
	MsgToRef     func(RefMsg) (Ref, error)
	MsgFromRef   func(Ref) RefMsg
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
)
