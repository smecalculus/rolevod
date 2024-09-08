package deal

import (
	"smecalculus/rolevod/app/seat"
)

type dealRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type dealRootData struct {
	ID       string        `db:"id"`
	Name     string        `db:"name"`
	Children []dealRefData `db:"-"`
	Prcs     []processData `db:"-"`
	Msgs     []messageData `db:"-"`
	Srvs     []serviceData `db:"-"`
}

type processData struct {
	ID string `db:"id"`
}

type messageData struct {
	ID string `db:"id"`
}

type serviceData struct {
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
	// goverter:ignore Seats Prcs Msgs Srvs
	DataToDealRoot func(dealRootData) (DealRoot, error)
	// goverter:ignore Prcs Msgs Srvs
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

type partRootData struct {
	Deal  dealRefData
	Seats []seat.SeatRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/app/seat:Data.*
var (
	DataToPartRoot   func(partRootData) (PartRoot, error)
	DataFromPartRoot func(PartRoot) partRootData
)
