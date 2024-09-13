package chnl

import (
	"smecalculus/rolevod/internal/state"
)

type SpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

type RootMsg struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	State *state.RefMsg `json:"state"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToSpec   func(SpecMsg) (Spec, error)
	MsgFromSpec func(Spec) SpecMsg
	MsgToRef    func(RefMsg) (Ref, error)
	MsgFromRef  func(Ref) *RefMsg
	// goverter:ignore PreID
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
)
