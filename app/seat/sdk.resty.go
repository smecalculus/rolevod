package seat

import (
	"fmt"
	"smecalculus/rolevod/lib/id"

	"github.com/go-resty/resty/v2"
)

// Adapter
type seatClient struct {
	resty *resty.Client
}

func newSeatClient() *seatClient {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &seatClient{r}
}

func NewSeatApi() SeatApi {
	return newSeatClient()
}

func (cl *seatClient) Create(spec SeatSpec) (SeatRoot, error) {
	req := MsgFromSeatSpec(spec)
	var res SeatRootMsg
	resp, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/seats")
	if err != nil {
		return SeatRoot{}, err
	}
	if resp.IsError() {
		return SeatRoot{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToSeatRoot(res)
}

func (c *seatClient) Retrieve(id id.ADT) (SeatRoot, error) {
	var res SeatRootMsg
	resp, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/seats/{id}")
	if err != nil {
		return SeatRoot{}, err
	}
	if resp.IsError() {
		return SeatRoot{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToSeatRoot(res)
}

func (c *seatClient) RetreiveAll() ([]SeatRef, error) {
	refs := []SeatRef{}
	return refs, nil
}

func (c *seatClient) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	resp, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/seats/{id}/kinships")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("received: %v", string(resp.Body()))
	}
	return nil
}
