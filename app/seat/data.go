package seat

type seatRootData struct {
	ID       string           `db:"id"`
	Name     string           `db:"name"`
	Children []seatTeaserData `db:"-"`
}

type seatTeaserData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type kinshipRootData struct {
	Parent   seatTeaserData
	Children []seatTeaserData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// seat
	DataToSeatTeaser    func(seatTeaserData) (SeatTeaser, error)
	DataFromSeatTeaser  func(SeatTeaser) seatTeaserData
	DataToSeatTeasers   func([]seatTeaserData) ([]SeatTeaser, error)
	DataFromSeatTeasers func([]SeatTeaser) []seatTeaserData
	// goverter:ignore Ctx Zc
	DataToSeatRoot   func(seatRootData) (SeatRoot, error)
	DataFromSeatRoot func(SeatRoot) seatRootData
	// kinship
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
