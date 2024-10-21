package sig

import (
	"fmt"
	"smecalculus/rolevod/lib/id"

	"github.com/go-resty/resty/v2"
)

// Adapter
type sigClient struct {
	resty *resty.Client
}

func newSigClient() *sigClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &sigClient{r}
}

func NewSigApi() Api {
	return newSigClient()
}

func (cl *sigClient) Create(spec Spec) (Root, error) {
	req := MsgFromSigSpec(spec)
	var res SigRootMsg
	resp, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/sigs")
	if err != nil {
		return Root{}, err
	}
	if resp.IsError() {
		return Root{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToSigRoot(res)
}

func (c *sigClient) Retrieve(id id.ADT) (Root, error) {
	var res SigRootMsg
	resp, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/sigs/{id}")
	if err != nil {
		return Root{}, err
	}
	if resp.IsError() {
		return Root{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToSigRoot(res)
}

func (c *sigClient) RetreiveAll() ([]Ref, error) {
	refs := []Ref{}
	return refs, nil
}

func (c *sigClient) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	resp, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/sigs/{id}/kinships")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("received: %v", string(resp.Body()))
	}
	return nil
}
