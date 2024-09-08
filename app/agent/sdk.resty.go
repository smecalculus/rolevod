package agent

import (
	"smecalculus/rolevod/lib/id"

	"github.com/go-resty/resty/v2"
)

// Adapter
type agentClient struct {
	resty *resty.Client
}

func newAgentClient() *agentClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &agentClient{r}
}

func NewAgentApi() AgentApi {
	return newAgentClient()
}

func (cl *agentClient) Create(spec AgentSpec) (AgentRoot, error) {
	req := MsgFromAgentSpec(spec)
	var res AgentRootMsg
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/agents")
	if err != nil {
		return AgentRoot{}, err
	}
	return MsgToAgentRoot(res)
}

func (c *agentClient) Retrieve(id id.ADT[ID]) (AgentRoot, error) {
	var res AgentRootMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/agents/{id}")
	if err != nil {
		return AgentRoot{}, err
	}
	return MsgToAgentRoot(res)
}

func (c *agentClient) RetreiveAll() ([]AgentRef, error) {
	refs := []AgentRef{}
	return refs, nil
}

func (c *agentClient) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	_, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/agents/{id}/kinships")
	return err
}
