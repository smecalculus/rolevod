package env

type RootData struct {
	ID   string
	Name string
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// goverter:ignore Tps Exps
	DataToRoot   func(RootData) (AR, error)
	DataFromRoot func(AR) RootData
)
