package alias

type rootData struct {
	Sym string
	ID  string
	Rev int64
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/rev:Convert.*
// goverter:extend smecalculus/rolevod/lib/sym:Convert.*
var (
	DataFromRoot func(Root) (rootData, error)
	DataToRoot   func(rootData) (Root, error)
)
