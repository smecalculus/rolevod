package deal_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/deal"
	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/seat"
)

var (
	roleApi = role.NewRoleApi()
	seatApi = seat.NewSeatApi()
	dealApi = deal.NewDealApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablishKinship(t *testing.T) {
	// given
	ps := deal.DealSpec{Name: "parent-deal"}
	pr, err := dealApi.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := deal.DealSpec{Name: "child-deal"}
	cr, err := dealApi.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := deal.KinshipSpec{
		ParentID:    pr.ID,
		ChildrenIDs: []id.ADT[deal.ID]{cr.ID},
	}
	err = dealApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := dealApi.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := deal.ToDealRef(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}

func TestTakeTransition(t *testing.T) {
	// given
	roleSpec := role.RoleSpec{
		Name:  "role-1",
		State: state.OneSpec{},
	}
	roleRoot, err := roleApi.Create(roleSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec := seat.SeatSpec{
		Name: "seat-1",
		Via: seat.ChanTp{
			Z:     "z-1",
			State: roleRoot.State,
		},
	}
	seatRoot, err := seatApi.Create(seatSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	dealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	dealRoot, err := dealApi.Create(dealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	partSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot.ID,
	}
	chnlRef, err := dealApi.Involve(partSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	tranSpec := deal.TranSpec{
		DealID: dealRoot.ID,
		Term: step.Wait{
			X:    chnlRef,
			Cont: nil,
		},
	}
	// when
	err = dealApi.Take(tranSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
}
