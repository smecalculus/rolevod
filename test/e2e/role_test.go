package e2e

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
	prs := role.RoleSpec{Name: "parent-role"}
	prr, err := roleApi.Create(prs)
	if err != nil {
		t.Fatal(err)
	}
	// and
	crs := role.RoleSpec{Name: "child-role", St: role.One{}}
	crr, err := roleApi.Create(crs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := role.KinshipSpec{
		Parent:   prr.ID,
		Children: []core.ID[role.Role]{crr.ID},
	}
	err = roleApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := roleApi.Retrieve(prr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := role.ToRoleTeaser(crr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", prr.Name, expectedChild, actual.Children)
	}
}
