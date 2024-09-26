package agent

type AgentSpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type AgentRootMsg struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Children []AgentRefMsg `json:"children"`
}

type AgentRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

type KinshipSpecMsg struct {
	ParentID    string   `param:"id" json:"parent"`
	ChildrenIDs []string `json:"children"`
}

type KinshipRootMsg struct {
	Parent   AgentRefMsg   `json:"parent"`
	Children []AgentRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
var (
	// agent
	MsgToAgentSpec    func(AgentSpecMsg) (AgentSpec, error)
	MsgFromAgentSpec  func(AgentSpec) AgentSpecMsg
	MsgToAgentRoot    func(AgentRootMsg) (AgentRoot, error)
	MsgFromAgentRoot  func(AgentRoot) AgentRootMsg
	MsgFromAgentRoots func([]AgentRoot) []AgentRootMsg
	// kinship
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
