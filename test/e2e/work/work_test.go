package work_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/work"
)

var (
	workApi = work.NewWorkApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablish(t *testing.T) {
	// given
	ps := work.WorkSpec{Name: "parent-work"}
	pr, err := workApi.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := work.WorkSpec{Name: "child-work"}
	cr, err := workApi.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := work.KinshipSpec{
		ParentID:    pr.ID,
		ChildrenIDs: []core.ID[work.Work]{cr.ID},
	}
	err = workApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := workApi.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := work.ToWorkTeaser(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}
