package agent

type agentRootData struct {
	ID       string         `db:"id"`
	Name     string         `db:"name"`
	Children []agentRefData `db:"-"`
}

type agentRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type kinshipRootData struct {
	Parent   agentRefData
	Children []agentRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
var (
	// agent
	DataToAgentRef    func(agentRefData) (AgentRef, error)
	DataFromAgentRef  func(AgentRef) agentRefData
	DataToAgentRefs   func([]agentRefData) ([]AgentRef, error)
	DataFromAgentRefs func([]AgentRef) []agentRefData
	DataToAgentRoot   func(agentRootData) (AgentRoot, error)
	DataFromAgentRoot func(AgentRoot) agentRootData
	// kinship
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
