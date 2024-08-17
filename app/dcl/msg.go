package dcl

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

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/app/dcl:to.*
var (
	ToSpec     func(SpecMsg) Spec
	ToSpecMsg  func(Spec) SpecMsg
	ToRoot     func(RootMsg) (TpDef, error)
	ToRootMsg  func(TpDef) RootMsg
	ToRootMsgs func([]TpDef) []RootMsg
)
