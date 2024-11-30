package role

import (
	"database/sql"
)

type refData struct {
	ID   string `db:"role_id"`
	Rev  int64  `db:"role_rev"`
	Name string `db:"role_name"`
}

type rootData struct {
	ID      string         `db:"role_id"`
	Rev     int64          `db:"role_rev"`
	Name    string         `db:"role_name"`
	StateID string         `db:"state_id"`
	WholeID sql.NullString `db:"whole_id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/rev:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
var (
	DataToRef     func(refData) (Ref, error)
	DataFromRef   func(Ref) (refData, error)
	DataToRefs    func([]refData) ([]Ref, error)
	DataFromRefs  func([]Ref) ([]refData, error)
	DataToRoot    func(rootData) (Root, error)
	DataFromRoot  func(Root) (rootData, error)
	DataToRoots   func([]rootData) ([]Root, error)
	DataFromRoots func([]Root) ([]rootData, error)
)
