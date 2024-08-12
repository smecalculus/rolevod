package env

type SpecMsg struct {
	Name string `json:"name"`
}

type RootMsg struct {
	Id string `json:"id"`
}

// goverter:converter
// goverter:output:package smecalculus/rolevod/app/env
// goverter:extend smecalculus/rolevod/app/env:to.*
type msgConverter interface {
	ToSpec(SpecMsg) Spec
	ToSpecMsg(Spec) SpecMsg
	ToRoot(RootMsg) (Root, error)
	ToRootMsg(Root) RootMsg
}
