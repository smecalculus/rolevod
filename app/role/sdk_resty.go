package role

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
		Post("/roles")
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
		Get("/roles/{id}")
	if err != nil {
		return Root{}, err
	}
	if resp.IsError() {
		return Root{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToRoot(res)
}

func (c *clientResty) Update(root Root) error {
	return nil
}

func (c *clientResty) RetreiveAll() ([]Ref, error) {
	rts := []Ref{}
	return rts, nil
}

func (c *clientResty) Establish(ks KinshipSpec) error {
	req := MsgFromKinshipSpec(ks)
	resp, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/roles/{id}/kinships")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("received: %v", string(resp.Body()))
	}
	return nil
}
