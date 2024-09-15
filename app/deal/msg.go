package deal

import (
	"smecalculus/rolevod/app/seat"
	"smecalculus/rolevod/internal/step"

	valid "github.com/go-ozzo/ozzo-validation"
)

type DealSpecMsg struct {
	Name  string            `json:"name"`
	Seats []seat.SeatRefMsg `json:"seats"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type DealRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
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
	Children []DealRefMsg      `json:"children"`
	Seats    []seat.SeatRefMsg `json:"seats"`
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
	ParentID    string   `param:"id" json:"parent"`
	ChildrenIDs []string `json:"children"`
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
	DealID string `param:"id" json:"deal_id"`
	SeatID string `json:"seat_id"`
}

type PartRootMsg struct {
	Deal DealRefMsg      `json:"deal"`
	Seat seat.SeatRefMsg `json:"seat"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/app/seat:To.*
// goverter:extend smecalculus/rolevod/app/seat:Msg.*
var (
	MsgFromPartSpec func(PartSpec) PartSpecMsg
	MsgToPartSpec   func(PartSpecMsg) (PartSpec, error)
	MsgFromPartRoot func(PartRoot) PartRootMsg
	MsgToPartRoot   func(PartRootMsg) (PartRoot, error)
)

type TransitionMsg struct {
	Deal DealRefMsg    `json:"deal"`
	Term *step.TermMsg `json:"term"`
}

func (v TransitionMsg) Validate() error {
	return valid.ValidateStruct(&v,
		valid.Field(&v.Deal, valid.Required),
		valid.Field(&v.Term, valid.NotNil),
	)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/internal/step:Msg.*
var (
	MsgFromTransition func(Transition) TransitionMsg
	MsgToTransition   func(TransitionMsg) (Transition, error)
)
