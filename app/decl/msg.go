package decl

type SpecMsg struct {
	Name string `json:"name"`
}

type RootMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetMsg struct {
	ID string `param:"id"`
}

// goverter:converter
// goverter:output:package smecalculus/rolevod/app/decl
// goverter:extend smecalculus/rolevod/app/decl:To.*
type MsgConverter interface {
	ToSpec(SpecMsg) Spec
	ToSpecMsg(Spec) SpecMsg
	ToRoot(RootMsg) (Root, error)
	ToRootMsg(Root) RootMsg
	ToRootMsgs([]Root) []RootMsg
}
