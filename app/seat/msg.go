package seat

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

type SeatSpecMsg struct {
	Name string         `json:"name"`
	Via  chnl.SpecMsg   `json:"via"`
	Ctx  []chnl.SpecMsg `json:"ctx"`
}

func (mto SeatSpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Name, core.NameRequired...),
		validation.Field(&mto.Via, validation.Required),
		validation.Field(&mto.Ctx, core.CtxOptional...),
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

type SeatRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

type SeatRootMsg struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Via      chnl.SpecMsg   `json:"via"`
	Ctx      []chnl.SpecMsg `json:"ctx"`
	Children []SeatRefMsg   `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/app/role:Msg.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToID          func(string) (id.ADT, error)
	MsgFromID        func(id.ADT) string
	MsgToSeatSpec    func(SeatSpecMsg) (SeatSpec, error)
	MsgFromSeatSpec  func(SeatSpec) SeatSpecMsg
	MsgToSeatRef     func(SeatRefMsg) (SeatRef, error)
	MsgFromSeatRef   func(SeatRef) SeatRefMsg
	MsgToSeatRoot    func(SeatRootMsg) (SeatRoot, error)
	MsgFromSeatRoot  func(SeatRoot) SeatRootMsg
	MsgFromSeatRoots func([]SeatRoot) []SeatRootMsg
)

type KinshipSpecMsg struct {
	ParentID string   `json:"parent_id" param:"id"`
	ChildIDs []string `json:"child_ids"`
}

type KinshipRootMsg struct {
	Parent   SeatRefMsg   `json:"parent"`
	Children []SeatRefMsg `json:"children"`
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
