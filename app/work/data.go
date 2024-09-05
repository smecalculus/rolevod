package work

type workRootData struct {
	ID       string           `db:"id"`
	Name     string           `db:"name"`
	Children []workTeaserData `db:"-"`
}

type workTeaserData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type kinshipRootData struct {
	Parent   workTeaserData
	Children []workTeaserData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// work
	DataToWorkTeaser    func(workTeaserData) (WorkTeaser, error)
	DataFromWorkTeaser  func(WorkTeaser) workTeaserData
	DataToWorkTeasers   func([]workTeaserData) ([]WorkTeaser, error)
	DataFromWorkTeasers func([]WorkTeaser) []workTeaserData
	DataToWorkRoot      func(workRootData) (WorkRoot, error)
	DataFromWorkRoot    func(WorkRoot) workRootData
	// kinship
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
