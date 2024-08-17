package dcl

type RootData struct {
	ID   string
	Name string
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/app/dcl:to.*
var (
	FromRootData func(RootData) (TpDef, error)
	ToRootData   func(TpDef) RootData
)
