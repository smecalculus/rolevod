package dcl

type RootData struct {
	ID   string
	Name string
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToTpDef   func(RootData) (TpRoot, error)
	DataFromTpDef func(TpRoot) RootData
)
