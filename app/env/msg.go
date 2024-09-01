package env

import (
	"smecalculus/rolevod/app/dcl"
)

type EnvSpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" json:"id"`
}

type EnvRootMsg struct {
	ID   string            `json:"id"`
	Name string            `json:"name"`
	Tps  []dcl.TpTeaserMsg `json:"tps"`
	Exps []dcl.ExpRootMsg  `json:"exps"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/app/dcl:Msg.*
var (
	MsgToEnvSpec    func(EnvSpecMsg) EnvSpec
	MsgFromEnvSpec  func(EnvSpec) EnvSpecMsg
	MsgFromEnvRoot  func(EnvRoot) EnvRootMsg
	MsgFromEnvRoots func([]EnvRoot) []EnvRootMsg
)
