package chnl

type SpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

type RootMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromSpec  func(Spec) SpecMsg
	MsgToRef     func(*RefMsg) (Ref, error)
	MsgFromRef   func(Ref) *RefMsg
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
)
