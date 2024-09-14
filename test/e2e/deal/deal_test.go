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
	rs := role.RoleSpec{
		Name:  "role-1",
		State: state.One{},
	}
	rr, err := roleApi.Create(rs)
	if err != nil {
		t.Fatal(err)
	}
	// and
	ss := seat.SeatSpec{
		Name: "seat-1",
		Via: seat.ChanTp{
			Z:    "z",
			Role: role.ToRoleRef(rr),
		},
	}
	sr, err := seatApi.Create(ss)
	if err != nil {
		t.Fatal(err)
	}
	// and
	ds := deal.DealSpec{
		Name: "big-deal",
	}
	dr, err := dealApi.Create(ds)
	if err != nil {
		t.Fatal(err)
	}
	// and
	ps := deal.PartSpec{
		DealID: dr.ID,
		SeatID: sr.ID,
	}
	a, err := dealApi.Involve(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	transition := deal.Transition{
		Deal: deal.ToDealRef(dr),
		Term: step.Wait{
			X:    a,
			Cont: step.Close{},
		},
	}
	// when
	err = dealApi.Take(transition)
	if err != nil {
		t.Fatal(err)
	}
	// then
}
