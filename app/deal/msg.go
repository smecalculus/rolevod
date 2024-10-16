package deal

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/seat"
)

type DealSpecMsg struct {
	Name string `json:"name"`
}

func (mto DealSpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Name, core.NameRequired...),
	)
}

type RefMsg struct {
	ID string `json:"id" param:"id"`
}

func (mto RefMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.ID, id.Required...),
	)
}

type DealRefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

func (mto DealRefMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.ID, id.Required...),
		validation.Field(&mto.Name, core.NameRequired...),
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
// goverter:extend smecalculus/rolevod/lib/id:String.*
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
	ParentID string   `json:"parent_id" param:"id"`
	ChildIDs []string `json:"child_ids"`
}

func (mto KinshipSpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.ParentID, id.Required...),
		validation.Field(&mto.ChildIDs,
			validation.Required, validation.Length(0, 10),
			validation.Each(id.Required...)),
	)
}

type KinshipRootMsg struct {
	Parent   DealRefMsg   `json:"parent"`
	Children []DealRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
var (
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)

type PartSpecMsg struct {
	DealID string        `json:"deal_id" param:"id"`
	SeatID string        `json:"seat_id"`
	OID    string        `json:"oid"`
	Ctx    []chnl.RefMsg `json:"ctx"`
}

func (mto PartSpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.DealID, id.Required...),
		validation.Field(&mto.SeatID, id.Required...),
		validation.Field(&mto.Ctx, core.CtxOptional...),
	)
}

type PartRootMsg struct {
	ID     string        `json:"id"`
	DealID string        `json:"deal_id"`
	SeatID string        `json:"seat_id"`
	PAK    string        `json:"pak"`
	CAK    string        `json:"cak"`
	PID    string        `json:"pid"`
	Ctx    []chnl.RefMsg `json:"ctx"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:String.*
// goverter:extend smecalculus/rolevod/internal/chnl:Msg.*
// goverter:extend smecalculus/rolevod/app/seat:Msg.*
var (
	MsgFromPartSpec func(PartSpec) PartSpecMsg
	MsgToPartSpec   func(PartSpecMsg) (PartSpec, error)
)

type TranSpecMsg struct {
	DID  string       `json:"did"`
	PID  string       `json:"pid"`
	Key  string       `json:"key"`
	Term step.TermMsg `json:"term"`
}

func (mto TranSpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.DID, id.Required...),
		validation.Field(&mto.PID, id.Required...),
		validation.Field(&mto.Key, ak.Required...),
		validation.Field(&mto.Term, validation.Required),
	)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:String.*
// goverter:extend smecalculus/rolevod/internal/step:Msg.*
var (
	MsgFromTranSpec func(TranSpec) TranSpecMsg
	MsgToTranSpec   func(TranSpecMsg) (TranSpec, error)
)
