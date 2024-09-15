package seat

import (
	"smecalculus/rolevod/internal/state"
	id "smecalculus/rolevod/lib/id"
)

type SeatSpecMsg struct {
	Name string      `json:"name"`
	Via  ChanTpMsg   `json:"via"`
	Ctx  []ChanTpMsg `json:"ctx"`
}

type ChanTpMsg struct {
	Z     string        `json:"z"`
	State *state.RefMsg `json:"state"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type SeatRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

type SeatRootMsg struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Via      ChanTpMsg    `json:"via"`
	Ctx      []ChanTpMsg  `json:"ctx"`
	Children []SeatRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/app/role:Msg.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToID          func(string) (id.ADT[ID], error)
	MsgFromID        func(id.ADT[ID]) string
	MsgToSeatSpec    func(SeatSpecMsg) (SeatSpec, error)
	MsgFromSeatSpec  func(SeatSpec) SeatSpecMsg
	MsgToSeatRef     func(SeatRefMsg) (SeatRef, error)
	MsgFromSeatRef   func(SeatRef) SeatRefMsg
	MsgToSeatRoot    func(SeatRootMsg) (SeatRoot, error)
	MsgFromSeatRoot  func(SeatRoot) SeatRootMsg
	MsgFromSeatRoots func([]SeatRoot) []SeatRootMsg
	MsgToChanTp      func(ChanTpMsg) (ChanTp, error)
	MsgFromChanTp    func(ChanTp) ChanTpMsg
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
// goverter:extend to.*
var (
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
