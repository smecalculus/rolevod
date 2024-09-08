package deal

type dealRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type dealRootData struct {
	ID       string        `db:"id"`
	Name     string        `db:"name"`
	Children []dealRefData `db:"-"`
	Procs    []procData    `db:"-"`
	Msgs     []msgData     `db:"-"`
	Srvs     []srvData     `db:"-"`
}

type procData struct {
	ID string `db:"id"`
}

type msgData struct {
	ID string `db:"id"`
}

type srvData struct {
	ID string `db:"id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToDealRef    func(dealRefData) (DealRef, error)
	DataFromDealRef  func(DealRef) dealRefData
	DataToDealRefs   func([]dealRefData) ([]DealRef, error)
	DataFromDealRefs func([]DealRef) []dealRefData
	DataToDealRoot   func(dealRootData) (DealRoot, error)
	DataFromDealRoot func(DealRoot) dealRootData
)

type kinshipRootData struct {
	Parent   dealRefData
	Children []dealRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)

// TODO: здесь не место
type seatRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type partRootData struct {
	Deal  dealRefData
	Seats []seatRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToPartRoot   func(partRootData) (PartRoot, error)
	DataFromPartRoot func(PartRoot) partRootData
)
