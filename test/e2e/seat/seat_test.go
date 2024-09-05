package seat_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/seat"
)

var (
	seatApi = seat.NewSeatApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablish(t *testing.T) {
	// given
	ps := seat.SeatSpec{Name: "parent-seat"}
	pr, err := seatApi.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := seat.SeatSpec{Name: "child-seat"}
	cr, err := seatApi.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := seat.KinshipSpec{
		ParentID:    pr.ID,
		ChildrenIDs: []core.ID[seat.Seat]{cr.ID},
	}
	err = seatApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := seatApi.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := seat.ToSeatTeaser(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}
