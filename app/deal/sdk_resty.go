package deal

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/lib/id"
)

func NewDealApi() DealApi {
	return newDealClient()
}

// Adapter
type dealClient struct {
	resty *resty.Client
}

func newDealClient() *dealClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &dealClient{r}
}

func (cl *dealClient) Create(spec DealSpec) (DealRoot, error) {
	req := MsgFromDealSpec(spec)
	var res DealRootMsg
	resp, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/deals")
	if err != nil {
		return DealRoot{}, err
	}
	if resp.IsError() {
		return DealRoot{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToDealRoot(res)
}

func (c *dealClient) Retrieve(id id.ADT) (DealRoot, error) {
	var res DealRootMsg
	resp, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/deals/{id}")
	if err != nil {
		return DealRoot{}, err
	}
	if resp.IsError() {
		return DealRoot{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToDealRoot(res)
}

func (c *dealClient) RetreiveAll() ([]DealRef, error) {
	refs := []DealRef{}
	return refs, nil
}

func (c *dealClient) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	resp, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/deals/{id}/kinships")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("received: %v", string(resp.Body()))
	}
	return nil
}

func (c *dealClient) Involve(spec PartSpec) (chnl.Root, error) {
	req := MsgFromPartSpec(spec)
	var res chnl.RootMsg
	resp, err := c.resty.R().
		SetResult(&res).
		SetBody(&req).
		SetPathParam("id", req.Deal).
		Post("/deals/{id}/parts")
	if err != nil {
		return chnl.Root{}, err
	}
	if resp.IsError() {
		return chnl.Root{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return chnl.MsgToRoot(res)
}

func (c *dealClient) Take(spec TranSpec) error {
	req := MsgFromTranSpec(spec)
	resp, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.DID).
		Post("/deals/{id}/steps")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("received: %v", string(resp.Body()))
	}
	return nil
}
