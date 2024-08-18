package env

import (
	"smecalculus/rolevod/app/dcl"
)

type SpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" json:"id"`
}

type RootMsg struct {
	ID   string           `json:"id"`
	Name string           `json:"name"`
	Tps  []dcl.TpRootMsg  `json:"tps"`
	Exps []dcl.ExpRootMsg `json:"exps"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/app/dcl:Msg.*
var (
	MsgToSpec    func(SpecMsg) AS
	MsgFromSpec  func(AS) SpecMsg
	MsgFromRoot  func(AR) RootMsg
	MsgFromRoots func([]AR) []RootMsg
)
