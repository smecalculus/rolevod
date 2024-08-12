package env

type RootData struct {
	Id string
}

// goverter:converter
// goverter:output:file ./generated.go
// goverter:output:package smecalculus/rolevod/app/env
// goverter:extend smecalculus/rolevod/app/env:to.*
type dataConverter interface {
	ToRoot(RootData) (Root, error)
	ToRootData(Root) RootData
}
