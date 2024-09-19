package seat

import (
	"database/sql"
)

type SeatRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type seatRootData struct {
	ID       string         `db:"id"`
	Name     string         `db:"name"`
	Via      sql.NullString `db:"via"`
	Ctx      sql.NullString `db:"ctx"`
	Children []SeatRefData  `db:"-"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
// goverter:extend smecalculus/rolevod/internal/state:Json.*
// goverter:extend smecalculus/rolevod/internal/chnl:Json.*
var (
	DataToSeatRef    func(SeatRefData) (SeatRef, error)
	DataFromSeatRef  func(SeatRef) SeatRefData
	DataToSeatRefs   func([]SeatRefData) ([]SeatRef, error)
	DataFromSeatRefs func([]SeatRef) []SeatRefData
	DataToSeatRoot   func(seatRootData) (SeatRoot, error)
	DataFromSeatRoot func(SeatRoot) (seatRootData, error)
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
