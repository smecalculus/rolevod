package role

import (
	"smecalculus/rolevod/lib/id"

	"github.com/go-resty/resty/v2"
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

func (cl *roleClient) Create(rs RoleSpec) (RoleRoot, error) {
	req := MsgFromRoleSpec(rs)
	var res RoleRootMsg
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/roles")
	if err != nil {
		return RoleRoot{}, err
	}
	return MsgToRoleRoot(res)
}

func (c *roleClient) Retrieve(id id.ADT[ID]) (RoleRoot, error) {
	var res RoleRootMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/roles/{id}")
	if err != nil {
		return RoleRoot{}, err
	}
	return MsgToRoleRoot(res)
}

func (c *roleClient) Update(rr RoleRoot) error {
	return nil
}

func (c *roleClient) RetreiveAll() ([]RoleRef, error) {
	rts := []RoleRef{}
	return rts, nil
}

func (c *roleClient) Establish(ks KinshipSpec) error {
	req := MsgFromKinshipSpec(ks)
	_, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.Parent).
		Post("/roles/{id}/kinships")
	return err
}
