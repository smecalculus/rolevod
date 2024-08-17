package env

type RootData struct {
	ID   string
	Name string
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/app/env:to.*
var (
	// goverter:ignore Decls
	FromRootData func(RootData) (Root, error)
	ToRootData   func(Root) RootData
)
