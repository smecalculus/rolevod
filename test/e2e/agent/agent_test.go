package agent_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/agent"
)

var (
	api = agent.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablish(t *testing.T) {
	// given
	ps := agent.Spec{Name: "parent-agent"}
	pr, err := api.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := agent.Spec{Name: "child-agent"}
	cr, err := api.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := agent.KinshipSpec{
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
	expectedChild := agent.ToAgentRef(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}
