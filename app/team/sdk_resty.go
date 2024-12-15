package team

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
		Post("/teams")
	if err != nil {
		return Root{}, err
	}
	return MsgToRoot(res)
}

func (c *clientResty) Retrieve(rid id.ADT) (Snap, error) {
	var res SnapMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", rid.String()).
		Get("/teams/{id}")
	if err != nil {
		return Snap{}, err
	}
	return MsgToSnap(res)
}

func (c *clientResty) RetreiveAll() ([]Ref, error) {
	refs := []Ref{}
	return refs, nil
}
