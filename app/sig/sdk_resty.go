package sig

import (
	"fmt"

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
	return MsgToRoot(res)
}

func (c *clientResty) Retrieve(id id.ADT) (Root, error) {
	var res RootMsg
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
	return MsgToRoot(res)
}

func (c *clientResty) RetreiveAll() ([]Ref, error) {
	refs := []Ref{}
	return refs, nil
}

func (c *clientResty) Establish(spec KinshipSpec) error {
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
