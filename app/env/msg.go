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

type IntroMsg struct {
	EnvID string `param:"id"`
	TpID  string `json:"tp_id"`
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
// goverter:extend smecalculus/rolevod/app/dcl:To.*
// goverter:extend smecalculus/rolevod/app/dcl:Msg.*
var (
	// env
	MsgToEnvSpec    func(EnvSpecMsg) EnvSpec
	MsgFromEnvSpec  func(EnvSpec) EnvSpecMsg
	MsgToEnvRoot    func(EnvRootMsg) (EnvRoot, error)
	MsgFromEnvRoot  func(EnvRoot) EnvRootMsg
	MsgFromEnvRoots func([]EnvRoot) []EnvRootMsg
	// intro
	MsgToIntro   func(IntroMsg) (TpIntro, error)
	MsgFromIntro func(TpIntro) IntroMsg
)
