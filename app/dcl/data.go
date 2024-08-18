package dcl

type RootData struct {
	ID   string
	Name string
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// goverter:ignore St
	DataToTpRoot   func(RootData) (TpRoot, error)
	DataFromTpRoot func(TpRoot) RootData
)
