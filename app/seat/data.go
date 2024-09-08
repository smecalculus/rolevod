package seat

type seatRootData struct {
	ID       string        `db:"id"`
	Name     string        `db:"name"`
	Children []seatRefData `db:"-"`
}

type seatRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type kinshipRootData struct {
	Parent   seatRefData
	Children []seatRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// seat
	DataToSeatRef    func(seatRefData) (SeatRef, error)
	DataFromSeatRef  func(SeatRef) seatRefData
	DataToSeatRefs   func([]seatRefData) ([]SeatRef, error)
	DataFromSeatRefs func([]SeatRef) []seatRefData
	// goverter:ignore Ctx Zc
	DataToSeatRoot   func(seatRootData) (SeatRoot, error)
	DataFromSeatRoot func(SeatRoot) seatRootData
	// kinship
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
