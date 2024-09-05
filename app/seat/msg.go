package seat

type SeatSpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type SeatRootMsg struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Children []SeatTeaserMsg `json:"children"`
}

type SeatTeaserMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `json:"name"`
}

type KinshipSpecMsg struct {
	ParentID    string   `param:"id" json:"parent"`
	ChildrenIDs []string `json:"children"`
}

type KinshipRootMsg struct {
	Parent   SeatTeaserMsg   `json:"parent"`
	Children []SeatTeaserMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// seat
	MsgToSeatSpec   func(SeatSpecMsg) (SeatSpec, error)
	MsgFromSeatSpec func(SeatSpec) SeatSpecMsg
	// goverter:ignore Ctx Zc
	MsgToSeatRoot    func(SeatRootMsg) (SeatRoot, error)
	MsgFromSeatRoot  func(SeatRoot) SeatRootMsg
	MsgFromSeatRoots func([]SeatRoot) []SeatRootMsg
	// kinship
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
