package force

import (
	"github.com/go-resty/resty/v2"

	"smecalculus/rolevod/lib/core"
)

// Adapter
type forceClient struct {
	resty *resty.Client
}

func newForceClient() *forceClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &forceClient{r}
}

func NewForceApi() ForceApi {
	return newForceClient()
}

func (cl *forceClient) Create(spec ForceSpec) (ForceRoot, error) {
	req := MsgFromForceSpec(spec)
	var res ForceRootMsg
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/forces")
	if err != nil {
		return ForceRoot{}, err
	}
	return MsgToForceRoot(res)
}

func (c *forceClient) Retrieve(id core.ID[Force]) (ForceRoot, error) {
	var res ForceRootMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/forces/{id}")
	if err != nil {
		return ForceRoot{}, err
	}
	return MsgToForceRoot(res)
}

func (c *forceClient) RetreiveAll() ([]ForceTeaser, error) {
	teasers := []ForceTeaser{}
	return teasers, nil
}

func (c *forceClient) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	_, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/forces/{id}/kinships")
	return err
}
