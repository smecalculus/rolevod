package force_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/force"
)

var (
	forceApi = force.NewForceApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablish(t *testing.T) {
	// given
	ps := force.ForceSpec{Name: "parent-force"}
	pr, err := forceApi.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := force.ForceSpec{Name: "child-force"}
	cr, err := forceApi.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := force.KinshipSpec{
		ParentID:    pr.ID,
		ChildrenIDs: []core.ID[force.Force]{cr.ID},
	}
	err = forceApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := forceApi.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := force.ToForceTeaser(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}
