package seat

import (
	"smecalculus/rolevod/internal/chnl"
)

type seatRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type seatRootData struct {
	ID       string          `db:"id"`
	Name     string          `db:"name"`
	PE       chnl.SpecData   `db:"pe"`
	CEs      []chnl.SpecData `db:"ces"`
	Children []seatRefData   `db:"-"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
var (
	DataToSeatRef     func(seatRefData) (SeatRef, error)
	DataFromSeatRef   func(SeatRef) seatRefData
	DataToSeatRefs    func([]seatRefData) ([]SeatRef, error)
	DataFromSeatRefs  func([]SeatRef) []seatRefData
	DataToSeatRoot    func(seatRootData) (SeatRoot, error)
	DataFromSeatRoot  func(SeatRoot) (seatRootData, error)
	DataToSeatRoots   func([]seatRootData) ([]SeatRoot, error)
	DataFromSeatRoots func([]SeatRoot) ([]seatRootData, error)
)

type kinshipRootData struct {
	Parent   seatRefData
	Children []seatRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
