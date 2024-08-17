package env

import (
	"smecalculus/rolevod/app/dcl"
)

type SpecMsg struct {
	Name string `json:"name"`
}

type RootMsg struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	Decls []dcl.RootMsg `json:"decls"`
}

type GetMsg struct {
	ID string `param:"id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/app/env:to.*
// goverter:extend smecalculus/rolevod/app/dcl:To.*
var (
	ToSpec     func(SpecMsg) Spec
	ToSpecMsg  func(Spec) SpecMsg
	ToRoot     func(RootMsg) (Root, error)
	ToRootMsg  func(Root) RootMsg
	ToRootMsgs func([]Root) []RootMsg
)
