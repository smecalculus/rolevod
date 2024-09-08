package chnl

type RootData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type RefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToRef    func(RefData) (Ref, error)
	DataFromRef  func(Ref) RefData
	DataToRefs   func([]RefData) ([]Ref, error)
	DataFromRefs func([]Ref) []RefData
	DataToRoot   func(RootData) (Root, error)
	DataFromRoot func(Root) RootData
)
