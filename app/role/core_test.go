package role

import (
	"log/slog"
	"testing"

	"smecalculus/rolevod/lib/core"
)

func TestKinshipEstalish(t *testing.T) {
	s := newRoleService(&roleRepoFixture{}, &kinshipRepoFixture{}, slog.Default())
	children := []core.ID[Role]{core.New[Role]()}
	s.Establish(KinshipSpec{Parent: core.New[Role](), Children: children})
}

type roleRepoFixture struct {
}

func (r *roleRepoFixture) Insert(rr RoleRoot) error {
	return nil
}
func (r *roleRepoFixture) SelectById(id core.ID[Role]) (RoleRoot, error) {
	return RoleRoot{}, nil
}
func (r *roleRepoFixture) SelectChildren(id core.ID[Role]) ([]RoleTeaser, error) {
	return []RoleTeaser{}, nil
}
func (r *roleRepoFixture) SelectAll() ([]RoleTeaser, error) {
	return []RoleTeaser{}, nil
}

type kinshipRepoFixture struct {
}

func (r *kinshipRepoFixture) Insert(kr KinshipRoot) error {
	return nil
}
