package decl

type RootData struct {
	ID   string
	Name string
}

// goverter:converter
// goverter:output:package smecalculus/rolevod/app/decl
// goverter:extend smecalculus/rolevod/app/decl:To.*
type dataConverter interface {
	ToRoot(RootData) (Root, error)
	ToRootData(Root) RootData
}
