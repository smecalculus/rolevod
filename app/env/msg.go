package env

import (
	"smecalculus/rolevod/app/decl"
)

type SpecMsg struct {
	Name string `json:"name"`
}

type RootMsg struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Decls []decl.RootMsg `json:"decls"`
}

type GetMsg struct {
	ID string `param:"id"`
}

// goverter:converter
// goverter:output:package smecalculus/rolevod/app/env
// goverter:extend smecalculus/rolevod/app/env:to.*
// goverter:extend smecalculus/rolevod/app/decl:To.*
type MsgConverter interface {
	ToSpec(SpecMsg) Spec
	ToSpecMsg(Spec) SpecMsg
	ToRoot(RootMsg) (Root, error)
	ToRootMsg(Root) RootMsg
	ToRootMsgs([]Root) []RootMsg
	// see https://github.com/jmattheis/goverter/issues/159
	ToDeclRoot(decl.RootMsg) (decl.TpDef, error)
	ToDeclRootMsg(decl.TpDef) decl.RootMsg
	ToDeclRootMsgs([]decl.TpDef) []decl.RootMsg
}
