package agent

import (
	"github.com/go-resty/resty/v2"

	"smecalculus/rolevod/lib/id"
)

// Adapter
type clientResty struct {
	resty *resty.Client
}

func newClientResty() *clientResty {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &clientResty{r}
}

func NewAPI() API {
	return newClientResty()
}

func (cl *clientResty) Create(spec Spec) (Root, error) {
	req := MsgFromSpec(spec)
	var res RootMsg
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/agents")
	if err != nil {
		return Root{}, err
	}
	return MsgToRoot(res)
}

func (c *clientResty) Retrieve(id id.ADT) (Root, error) {
	var res RootMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/agents/{id}")
	if err != nil {
		return Root{}, err
	}
	return MsgToRoot(res)
}

func (c *clientResty) RetreiveAll() ([]Ref, error) {
	refs := []Ref{}
	return refs, nil
}

func (c *clientResty) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	_, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/agents/{id}/kinships")
	return err
}
