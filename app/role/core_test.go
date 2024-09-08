package role

import (
	"log/slog"
	"testing"

	"smecalculus/rolevod/lib/id"
)

func TestKinshipEstalish(t *testing.T) {
	s := newRoleService(&roleRepoFixture{}, &kinshipRepoFixture{}, slog.Default())
	children := []id.ADT[ID]{id.New[ID]()}
	s.Establish(KinshipSpec{Parent: id.New[ID](), Children: children})
}

type roleRepoFixture struct {
}

func (r *roleRepoFixture) Insert(rr RoleRoot) error {
	return nil
}
func (r *roleRepoFixture) SelectById(id id.ADT[ID]) (RoleRoot, error) {
	return RoleRoot{}, nil
}
func (r *roleRepoFixture) SelectChildren(id id.ADT[ID]) ([]RoleRef, error) {
	return []RoleRef{}, nil
}
func (r *roleRepoFixture) SelectAll() ([]RoleRef, error) {
	return []RoleRef{}, nil
}

type kinshipRepoFixture struct {
}

func (r *kinshipRepoFixture) Insert(kr KinshipRoot) error {
	return nil
}
