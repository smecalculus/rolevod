package env

type RootData struct {
	ID   string
	Name string
}

// goverter:converter
// goverter:output:package smecalculus/rolevod/app/env
// goverter:extend smecalculus/rolevod/app/env:to.*
type dataConverter interface {
	// goverter:ignore Decls
	ToRoot(RootData) (Root, error)
	ToRootData(Root) RootData
}
