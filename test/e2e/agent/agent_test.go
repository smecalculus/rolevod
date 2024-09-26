package agent_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/agent"
)

var (
	agentApi = agent.NewAgentApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablish(t *testing.T) {
	// given
	ps := agent.AgentSpec{Name: "parent-agent"}
	pr, err := agentApi.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := agent.AgentSpec{Name: "child-agent"}
	cr, err := agentApi.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := agent.KinshipSpec{
		ParentID:    pr.ID,
		ChildrenIDs: []id.ADT{cr.ID},
	}
	err = agentApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := agentApi.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := agent.ToAgentRef(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}
