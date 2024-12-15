package crew

type rootData struct {
	ID       string    `db:"id"`
	Name     string    `db:"name"`
	Children []refData `db:"-"`
}

type refData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type kinshipRootData struct {
	Parent   refData
	Children []refData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	// crew
	DataToRef    func(refData) (Ref, error)
	DataFromRef  func(Ref) refData
	DataToRefs   func([]refData) ([]Ref, error)
	DataFromRefs func([]Ref) []refData
	DataToRoot   func(rootData) (Root, error)
	DataFromRoot func(Root) rootData
	// kinship
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
