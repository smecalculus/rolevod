package pool

import (
	"database/sql"

	"smecalculus/rolevod/internal/proc"
)

type refData struct {
	PoolID string `json:"pool_id"`
	Title  string `json:"title"`
}

type subSnapData struct {
	PoolID string    `db:"pool_id"`
	Title  string    `db:"title"`
	Subs   []refData `db:"subs"`
}

type assetSnapData struct {
	PoolID string `db:"pool_id"`
	Title  string `db:"title"`
}

type rootData struct {
	PoolID string         `db:"pool_id"`
	Revs   []int          `db:"revs"`
	Title  string         `db:"title"`
	SupID  sql.NullString `db:"sup_pool_id"`
}

type assetModData struct {
	InPoolID  string `db:"pool_id"`
	OutPoolID string `db:"ex_pool_id"`
	Rev       int    `db:"rev"`
	EPs       []epData
}

type epData struct {
	ProcID   string `db:"proc_id"`
	ProcPH   string `db:"proc_ph"`
	ChnlID   string `db:"chnl_id"`
	StateID  string `db:"state_id"`
	SrvID    string `db:"srv_id"`
	SrvRevs  []int  `db:"srv_revs"`
	ClntID   string `db:"clnt_id"`
	ClntRevs []int  `db:"clnt_revs"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	DataToRef        func(refData) (Ref, error)
	DataFromRef      func(Ref) refData
	DataToRefs       func([]refData) ([]Ref, error)
	DataFromRefs     func([]Ref) []refData
	DataToRoot       func(rootData) (Root, error)
	DataFromRoot     func(Root) rootData
	DataToSubSnap    func(subSnapData) (SubSnap, error)
	DataFromSubSnap  func(SubSnap) subSnapData
	DataFromAssetMod func(AssetMod) assetModData
	DataToAssetSnap  func(assetSnapData) (AssetSnap, error)
	DataToEPs        func([]epData) ([]proc.EP, error)
)
