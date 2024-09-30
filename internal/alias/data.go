package alias

type rootData struct {
	Sym string
	ID  string
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/sym:String.*
var (
	DataFromRoot func(Root) (rootData, error)
	DataToRoot   func(rootData) (Root, error)
)
