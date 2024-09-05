package work

import (
	"github.com/go-resty/resty/v2"

	"smecalculus/rolevod/lib/core"
)

// Adapter
type workClient struct {
	resty *resty.Client
}

func newWorkClient() *workClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &workClient{r}
}

func NewWorkApi() WorkApi {
	return newWorkClient()
}

func (cl *workClient) Create(spec WorkSpec) (WorkRoot, error) {
	req := MsgFromWorkSpec(spec)
	var res WorkRootMsg
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/works")
	if err != nil {
		return WorkRoot{}, err
	}
	return MsgToWorkRoot(res)
}

func (c *workClient) Retrieve(id core.ID[Work]) (WorkRoot, error) {
	var res WorkRootMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/works/{id}")
	if err != nil {
		return WorkRoot{}, err
	}
	return MsgToWorkRoot(res)
}

func (c *workClient) RetreiveAll() ([]WorkTeaser, error) {
	teasers := []WorkTeaser{}
	return teasers, nil
}

func (c *workClient) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	_, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/works/{id}/kinships")
	return err
}
