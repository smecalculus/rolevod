package sig_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/app/sig"
	"smecalculus/rolevod/lib/id"
)

var (
	sigApi = sig.NewSigApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablish(t *testing.T) {
	// given
	ps := sig.Spec{FQN: "parent-sig"}
	pr, err := sigApi.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := sig.Spec{FQN: "child-sig"}
	cr, err := sigApi.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := sig.KinshipSpec{
		ParentID: pr.ID,
		ChildIDs: []id.ADT{cr.ID},
	}
	err = sigApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := sigApi.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := sig.ConvertRootToRef(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}
