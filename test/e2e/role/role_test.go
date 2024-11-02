package role_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"

	"smecalculus/rolevod/app/role"
)

var (
	api = role.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestKinshipEstablishment(t *testing.T) {
	// given
	parSpec := role.Spec{FQN: "parent-role"}
	parRoot, err := api.Create(parSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	childSpec := role.Spec{FQN: "child-role", State: state.OneSpec{}}
	childRoot, err := api.Create(childSpec)
	if err != nil {
		t.Fatal(err)
	}
	// when
	kinshipSpec := role.KinshipSpec{
		ParentID: parRoot.ID,
		ChildIDs: []id.ADT{childRoot.ID},
	}
	err = api.Establish(kinshipSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := api.Retrieve(parRoot.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := role.ConverRootToRef(childRoot)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", parRoot.Name, expectedChild, actual.Children)
	}
}
