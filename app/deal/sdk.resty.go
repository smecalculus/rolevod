package deal

import (
	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/lib/id"

	"github.com/go-resty/resty/v2"
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
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/deals")
	if err != nil {
		return DealRoot{}, err
	}
	return MsgToDealRoot(res)
}

func (c *dealClient) Retrieve(id id.ADT[ID]) (DealRoot, error) {
	var res DealRootMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/deals/{id}")
	if err != nil {
		return DealRoot{}, err
	}
	return MsgToDealRoot(res)
}

func (c *dealClient) RetreiveAll() ([]DealRef, error) {
	refs := []DealRef{}
	return refs, nil
}

func (c *dealClient) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	_, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/deals/{id}/kinships")
	return err
}

func (c *dealClient) Involve(spec PartSpec) (chnl.Ref, error) {
	req := MsgFromPartSpec(spec)
	var res chnl.RefMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetBody(&req).
		SetPathParam("id", req.DealID).
		Post("/deals/{id}/parts")
	if err != nil {
		return chnl.Ref{}, err
	}
	return chnl.MsgToRef(res)
}

func (c *dealClient) Take(spec TranSpec) error {
	req := MsgFromTranSpec(spec)
	_, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.DealID).
		Post("/deals/{id}/steps")
	return err
}
