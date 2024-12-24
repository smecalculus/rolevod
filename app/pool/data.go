package pool

import (
	"database/sql"
)

type refData struct {
	ID    string `json:"pool_id"`
	Rev   int64  `json:"rev"`
	Title string `json:"title"`
}

type snapData struct {
	ID    string    `db:"pool_id"`
	Title string    `db:"title"`
	Subs  []refData `db:"subs"`
}

type rootData struct {
	ID    string         `db:"pool_id"`
	Rev   int64          `db:"rev"`
	Title string         `db:"title"`
	SupID sql.NullString `db:"sup_id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	DataToRef    func(refData) (Ref, error)
	DataFromRef  func(Ref) refData
	DataToRefs   func([]refData) ([]Ref, error)
	DataFromRefs func([]Ref) []refData
	DataToRoot   func(rootData) (Root, error)
	DataFromRoot func(Root) rootData
	DataToSnap   func(snapData) (Snap, error)
	DataFromSnap func(Snap) snapData
)
