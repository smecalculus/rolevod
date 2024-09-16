package role

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"smecalculus/rolevod/lib/id"
)

// Adapter
type roleClient struct {
	resty *resty.Client
}

func newRoleClient() *roleClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &roleClient{r}
}

func NewRoleApi() RoleApi {
	return newRoleClient()
}

func (cl *roleClient) Create(spec RoleSpec) (RoleRoot, error) {
	req := MsgFromRoleSpec(spec)
	var res RoleRootMsg
	resp, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/roles")
	if err != nil {
		return RoleRoot{}, err
	}
	if resp.IsError() {
		return RoleRoot{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToRoleRoot(res)
}

func (c *roleClient) Retrieve(id id.ADT[ID]) (RoleRoot, error) {
	var res RoleRootMsg
	resp, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/roles/{id}")
	if err != nil {
		return RoleRoot{}, err
	}
	if resp.IsError() {
		return RoleRoot{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToRoleRoot(res)
}

func (c *roleClient) Update(root RoleRoot) error {
	return nil
}

func (c *roleClient) RetreiveAll() ([]RoleRef, error) {
	rts := []RoleRef{}
	return rts, nil
}

func (c *roleClient) Establish(ks KinshipSpec) error {
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
