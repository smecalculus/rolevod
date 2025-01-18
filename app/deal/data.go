package deal

type refData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type rootData struct {
	ID       string    `db:"id"`
	Name     string    `db:"name"`
	Children []refData `db:"-"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	DataToRef    func(refData) (Ref, error)
	DataFromRef  func(Ref) refData
	DataToRefs   func([]refData) ([]Ref, error)
	DataFromRefs func([]Ref) []refData
	// goverter:ignore Sigs
	DataToRoot   func(rootData) (Root, error)
	DataFromRoot func(Root) rootData
)
