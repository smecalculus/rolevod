package deal_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/deal"
)

var (
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
	ds := deal.DealSpec{Name: "parent-deal"}
	dr, err := dealApi.Create(ds)
	if err != nil {
		t.Fatal(err)
	}
	// and
	x := chnl.Root{
		ID:    id.New[chnl.ID](),
		Name:  "x",
		State: state.One{},
	}
	// and
	tran := deal.Transition{
		Deal: deal.ToDealRef(dr),
		Term: step.Wait{
			X:    chnl.ToRef(x),
			Cont: step.Close{},
		},
	}
	// when
	err = dealApi.Take(tran)
	if err != nil {
		t.Fatal(err)
	}
	// then
}
