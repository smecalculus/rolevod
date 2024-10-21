package deal

type dealRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type dealRootData struct {
	ID       string        `db:"id"`
	Name     string        `db:"name"`
	Children []dealRefData `db:"-"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
var (
	DataToDealRef    func(dealRefData) (DealRef, error)
	DataFromDealRef  func(DealRef) dealRefData
	DataToDealRefs   func([]dealRefData) ([]DealRef, error)
	DataFromDealRefs func([]DealRef) []dealRefData
	// goverter:ignore Sigs
	DataToDealRoot   func(dealRootData) (DealRoot, error)
	DataFromDealRoot func(DealRoot) dealRootData
)

type kinshipRootData struct {
	Parent   dealRefData
	Children []dealRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
