package role_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/role"

	"smecalculus/rolevod/internal/state"
)

var (
	roleApi = role.NewRoleApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestKinshipEstablishment(t *testing.T) {
	// given
	parSpec := role.RoleSpec{Name: "parent-role"}
	parRoot, err := roleApi.Create(parSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	childSpec := role.RoleSpec{Name: "child-role", State: state.OneSpec{}}
	childRoot, err := roleApi.Create(childSpec)
	if err != nil {
		t.Fatal(err)
	}
	// when
	kinshipSpec := role.KinshipSpec{
		ParentID:    parRoot.ID,
		ChildrenIDs: []id.ADT[role.ID]{childRoot.ID},
	}
	err = roleApi.Establish(kinshipSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := roleApi.Retrieve(parRoot.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := role.ConverRootToRef(childRoot)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", parRoot.Name, expectedChild, actual.Children)
	}
}
