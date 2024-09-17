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

// type chanTpData struct {
// 	Z     string         `json:"z"`
// 	State *state.RefData `json:"state"`
// }

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
	// DataToChanTp     func(chanTpData) (ChanTp, error)
	// DataFromChanTp   func(ChanTp) (chanTpData, error)
	// DataToChanTps    func([]chanTpData) ([]ChanTp, error)
	// DataFromChanTps  func([]ChanTp) ([]chanTpData, error)
)

// func jsonFromChanTp(rel ChanTp) (string, error) {
// 	dto, err := DataFromChanTp(rel)
// 	if err != nil {
// 		return "", err
// 	}
// 	data, err := json.Marshal(dto)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(data), nil
// }

// func jsonToChanTp(data string) (ChanTp, error) {
// 	var dto chanTpData
// 	err := json.Unmarshal([]byte(data), &dto)
// 	if err != nil {
// 		return ChanTp{}, err
// 	}
// 	return DataToChanTp(dto)
// }

// func jsonFromChanTps(rels []ChanTp) (sql.NullString, error) {
// 	null := sql.NullString{}
// 	if len(rels) == 0 {
// 		return null, nil
// 	}
// 	dtos, err := DataFromChanTps(rels)
// 	if err != nil {
// 		return null, err
// 	}
// 	jsn, err := json.Marshal(dtos)
// 	if err != nil {
// 		return null, err
// 	}
// 	return sql.NullString{String: string(jsn), Valid: true}, nil
// }

// func jsonToChanTps(jsn sql.NullString) ([]ChanTp, error) {
// 	if !jsn.Valid {
// 		return []ChanTp{}, nil
// 	}
// 	var dtos []chanTpData
// 	err := json.Unmarshal([]byte(jsn.String), &dtos)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return DataToChanTps(dtos)
// }

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
