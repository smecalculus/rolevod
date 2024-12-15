package crew

type SpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type RootMsg struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Children []AgentRefMsg `json:"children"`
}

type AgentRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

type KinshipSpecMsg struct {
	ParentID string   `param:"id" json:"parent"`
	ChildIDs []string `json:"children"`
}

type KinshipRootMsg struct {
	Parent   AgentRefMsg   `json:"parent"`
	Children []AgentRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	// crew
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromSpec  func(Spec) SpecMsg
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
	// kinship
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
