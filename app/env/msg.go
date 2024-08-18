package env

import (
	"smecalculus/rolevod/app/dcl"
)

type SpecMsg struct {
	Name string `json:"name"`
}

type RootMsg struct {
	ID   string        `json:"id"`
	Name string        `json:"name"`
	Tps  []dcl.RootMsg `json:"tps"`
	Exps []dcl.RootMsg `json:"exps"`
}

type GetMsg struct {
	ID string `param:"id"`
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
