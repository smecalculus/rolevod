package seat

type seatRootData struct {
	ID       string        `db:"id"`
	Name     string        `db:"name"`
	Children []SeatRefData `db:"-"`
}

type SeatRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToSeatRef    func(SeatRefData) (SeatRef, error)
	DataFromSeatRef  func(SeatRef) SeatRefData
	DataToSeatRefs   func([]SeatRefData) ([]SeatRef, error)
	DataFromSeatRefs func([]SeatRef) []SeatRefData
	// goverter:ignore Ctx Zc
	DataToSeatRoot   func(seatRootData) (SeatRoot, error)
	DataFromSeatRoot func(SeatRoot) seatRootData
)

type kinshipRootData struct {
	Parent   SeatRefData
	Children []SeatRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
