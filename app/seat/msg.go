package seat

type SeatSpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type SeatRootMsg struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Children []SeatRefMsg `json:"children"`
}

type SeatRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	MsgToSeatSpec   func(SeatSpecMsg) (SeatSpec, error)
	MsgFromSeatSpec func(SeatSpec) SeatSpecMsg
	MsgToSeatRef    func(SeatRefMsg) (SeatRef, error)
	MsgFromSeatRef  func(SeatRef) SeatRefMsg
	// goverter:ignore Ctx Zc
	MsgToSeatRoot    func(SeatRootMsg) (SeatRoot, error)
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
// goverter:extend to.*
var (
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
