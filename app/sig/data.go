package sig

import (
	"smecalculus/rolevod/internal/chnl"
)

type sigRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type sigRootData struct {
	ID       string          `db:"id"`
	Name     string          `db:"name"`
	PE       chnl.SpecData   `db:"pe"`
	CEs      []chnl.SpecData `db:"ces"`
	Children []sigRefData    `db:"-"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
var (
	DataToSigRef     func(sigRefData) (Ref, error)
	DataFromSigRef   func(Ref) sigRefData
	DataToSigRefs    func([]sigRefData) ([]Ref, error)
	DataFromSigRefs  func([]Ref) []sigRefData
	DataToSigRoot    func(sigRootData) (Root, error)
	DataFromSigRoot  func(Root) (sigRootData, error)
	DataToSigRoots   func([]sigRootData) ([]Root, error)
	DataFromSigRoots func([]Root) ([]sigRootData, error)
)

type kinshipRootData struct {
	Parent   sigRefData
	Children []sigRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
