package chnl

import (
	"smecalculus/rolevod/lib/id"

	"github.com/go-resty/resty/v2"
)

// Adapter
type Client struct {
	resty *resty.Client
}

func newClient() *Client {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &Client{r}
}

func NewApi() Api {
	return newClient()
}

func (cl *Client) Create(spec Spec) (Root, error) {
	req := MsgFromSpec(spec)
	var res RootMsg
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/s")
	if err != nil {
		return Root{}, err
	}
	return MsgToRoot(res)
}

func (c *Client) Retrieve(id id.ADT[ID]) (Root, error) {
	var res RootMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/s/{id}")
	if err != nil {
		return Root{}, err
	}
	return MsgToRoot(res)
}

func (c *Client) RetreiveAll() ([]Ref, error) {
	refs := []Ref{}
	return refs, nil
}
