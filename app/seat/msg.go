package seat

import (
	valid "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

type SeatSpecMsg struct {
	Name string         `json:"name"`
	Via  chnl.SpecMsg   `json:"via"`
	Ctx  []chnl.SpecMsg `json:"ctx"`
}

func (mto *SeatSpecMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.Name, valid.Required, valid.Length(1, 64)),
		valid.Field(&mto.Via, valid.Required),
	)
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
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
	MsgToID         func(string) (id.ADT, error)
	MsgFromID       func(id.ADT) string
	// goverter:ignore Ctx
	MsgToSeatSpec   func(SeatSpecMsg) (SeatSpec, error)
	// goverter:ignore Ctx
	MsgFromSeatSpec func(SeatSpec) SeatSpecMsg
	MsgToSeatRef    func(SeatRefMsg) (SeatRef, error)
	MsgFromSeatRef  func(SeatRef) SeatRefMsg
	// goverter:ignore Ctx
	MsgToSeatRoot func(SeatRootMsg) (SeatRoot, error)
	// goverter:ignore Ctx
	MsgFromSeatRoot  func(SeatRoot) SeatRootMsg
	MsgFromSeatRoots func([]SeatRoot) []SeatRootMsg
)

type KinshipSpecMsg struct {
	ParentID    string   `param:"id" json:"parent"`
	ChildrenIDs []string `json:"children"`
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
