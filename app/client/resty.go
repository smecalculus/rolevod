package client

import (
	"github.com/go-resty/resty/v2"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/dcl"
	ws "smecalculus/rolevod/app/env"
)

// Adapter
type envClient struct {
	resty *resty.Client
}

func newEnvClient() *envClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &envClient{r}
}

func NewEnvApi() ws.EnvApi {
	return newEnvClient()
}

func (c *envClient) Create(spec ws.EnvSpec) (ws.EnvRoot, error) {
	req := ws.MsgFromEnvSpec(spec)
	var res ws.EnvRootMsg
	_, err := c.resty.R().
		SetBody(req).
		SetResult(&res).
		Post("/envs")
	if err != nil {
		return ws.EnvRoot{}, err
	}
	return ws.MsgToEnvRoot(res)
}

func (c *envClient) Retrieve(id core.ID[ws.AR]) (ws.EnvRoot, error) {
	var res ws.EnvRootMsg
	_, err := c.resty.R().
		SetPathParam("id", id.String()).
		SetResult(&res).
		Get("/envs/{id}")
	if err != nil {
		return ws.EnvRoot{}, err
	}
	return ws.MsgToEnvRoot(res)
}

func (c *envClient) RetreiveAll() ([]ws.EnvRoot, error) {
	roots := []ws.EnvRoot{}
	return roots, nil
}

func (c *envClient) Introduce(intro ws.TpIntro) error {
	req := ws.MsgFromIntro(intro)
	_, err := c.resty.R().
		SetPathParam("id", req.EnvID).
		SetBody(req).
		Post("/envs/{id}/intros")
	return err
}

// Adapter
type tpClient struct {
	resty *resty.Client
}

func newTpClient() *tpClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &tpClient{r}
}

func NewTpApi() dcl.TpApi {
	return newTpClient()
}

func (c *tpClient) Create(spec dcl.TpSpec) (dcl.TpRoot, error) {
	req := dcl.MsgFromTpSpec(spec)
	var res dcl.TpRootRaw
	_, err := c.resty.R().
		SetBody(req).
		SetResult(&res).
		Post("/tps")
	if err != nil {
		return dcl.TpRoot{}, err
	}
	return dcl.MsgToTpRoot(res)
}

func (c *tpClient) Update(root dcl.TpRoot) error {
	return nil
}

func (c *tpClient) Retrieve(id core.ID[dcl.AR]) (dcl.TpRoot, error) {
	root := dcl.TpRoot{}
	return root, nil
}

func (c *tpClient) RetreiveAll() ([]dcl.TpRoot, error) {
	roots := []dcl.TpRoot{}
	return roots, nil
}
