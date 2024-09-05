package role_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/role"
)

var (
	roleApi = role.NewRoleApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablish(t *testing.T) {
	// given
	ps := role.RoleSpec{Name: "parent-role"}
	pr, err := roleApi.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := role.RoleSpec{Name: "child-role", St: role.One{}}
	cr, err := roleApi.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := role.KinshipSpec{
		Parent:   pr.ID,
		Children: []core.ID[role.Role]{cr.ID},
	}
	err = roleApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := roleApi.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := role.ToRoleTeaser(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}