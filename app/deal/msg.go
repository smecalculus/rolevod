package deal

import (
	valid "github.com/go-ozzo/ozzo-validation"

	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/seat"
)

type DealSpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `json:"id" param:"id"`
}

type DealRefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

func (v DealRefMsg) Validate() error {
	return valid.ValidateStruct(&v,
		valid.Field(&v.ID, valid.Required, valid.Length(20, 20)),
		valid.Field(&v.Name, valid.Required, valid.Length(1, 64)),
	)
}

type DealRootMsg struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Seats    []seat.SeatRefMsg `json:"seats"`
	Children []DealRefMsg      `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/app/seat:Msg.*
var (
	MsgToDealSpec    func(DealSpecMsg) (DealSpec, error)
	MsgFromDealSpec  func(DealSpec) DealSpecMsg
	MsgToDealRef     func(DealRefMsg) (DealRef, error)
	MsgFromDealRef   func(DealRef) *DealRefMsg
	MsgToDealRoot    func(DealRootMsg) (DealRoot, error)
	MsgFromDealRoot  func(DealRoot) DealRootMsg
	MsgFromDealRoots func([]DealRoot) []DealRootMsg
)

type KinshipSpecMsg struct {
	ParentID    string   `json:"parent_id" param:"id"`
	ChildrenIDs []string `json:"children_ids"`
}

type KinshipRootMsg struct {
	Parent   DealRefMsg   `json:"parent"`
	Children []DealRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)

type PartSpecMsg struct {
	DealID string `json:"deal_id" param:"id"`
	SeatID string `json:"seat_id"`
}

type PartRootMsg struct {
	Deal DealRefMsg      `json:"deal"`
	Seat seat.SeatRefMsg `json:"seat"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/app/seat:Msg.*
var (
	MsgFromPartSpec func(PartSpec) PartSpecMsg
	MsgToPartSpec   func(PartSpecMsg) (PartSpec, error)
	MsgFromPartRoot func(PartRoot) PartRootMsg
	MsgToPartRoot   func(PartRootMsg) (PartRoot, error)
)

type TranSpecMsg struct {
	DealID string        `json:"deal_id"`
	Term   *step.TermMsg `json:"term"`
}

func (mto *TranSpecMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.DealID, valid.Required),
		valid.Field(&mto.Term, valid.Required),
	)
}

type TranRootMsg struct {
	Deal DealRefMsg    `json:"deal"`
	Term *step.TermMsg `json:"term"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/internal/step:Msg.*
var (
	MsgFromTranSpec func(TranSpec) TranSpecMsg
	MsgToTranSpec   func(TranSpecMsg) (TranSpec, error)
	MsgFromTranRoot func(TranRoot) TranRootMsg
	MsgToTranRoot   func(TranRootMsg) (TranRoot, error)
)
