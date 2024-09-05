package seat

import (
	"github.com/go-resty/resty/v2"

	"smecalculus/rolevod/lib/core"
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
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/seats")
	if err != nil {
		return SeatRoot{}, err
	}
	return MsgToSeatRoot(res)
}

func (c *seatClient) Retrieve(id core.ID[Seat]) (SeatRoot, error) {
	var res SeatRootMsg
	_, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/seats/{id}")
	if err != nil {
		return SeatRoot{}, err
	}
	return MsgToSeatRoot(res)
}

func (c *seatClient) RetreiveAll() ([]SeatTeaser, error) {
	teasers := []SeatTeaser{}
	return teasers, nil
}

func (c *seatClient) Establish(spec KinshipSpec) error {
	req := MsgFromKinshipSpec(spec)
	_, err := c.resty.R().
		SetBody(&req).
		SetPathParam("id", req.ParentID).
		Post("/seats/{id}/kinships")
	return err
}
