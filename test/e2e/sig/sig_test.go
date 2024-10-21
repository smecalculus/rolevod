package sig_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/sig"
)

var (
	api = sig.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablish(t *testing.T) {
	// given
	ps := sig.Spec{FQN: "parent-sig"}
	pr, err := api.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := sig.Spec{FQN: "child-sig"}
	cr, err := api.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := sig.KinshipSpec{
		ParentID: pr.ID,
		ChildIDs: []id.ADT{cr.ID},
	}
	err = api.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := api.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := sig.ConvertRootToRef(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}
