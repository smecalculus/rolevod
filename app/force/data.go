package force

type forceRootData struct {
	ID       string            `db:"id"`
	Name     string            `db:"name"`
	Children []forceTeaserData `db:"-"`
}

type forceTeaserData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type kinshipRootData struct {
	Parent   forceTeaserData
	Children []forceTeaserData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// force
	DataToForceTeaser    func(forceTeaserData) (ForceTeaser, error)
	DataFromForceTeaser  func(ForceTeaser) forceTeaserData
	DataToForceTeasers   func([]forceTeaserData) ([]ForceTeaser, error)
	DataFromForceTeasers func([]ForceTeaser) []forceTeaserData
	DataToForceRoot      func(forceRootData) (ForceRoot, error)
	DataFromForceRoot    func(ForceRoot) forceRootData
	// kinship
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
