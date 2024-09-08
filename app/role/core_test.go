package role

import (
	"log/slog"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

func TestKinshipEstalish(t *testing.T) {
	s := newRoleService(&roleRepoFixture{}, &stateRepoFixture{}, &kinshipRepoFixture{}, slog.Default())
	children := []id.ADT[ID]{id.New[ID]()}
	s.Establish(KinshipSpec{Parent: id.New[ID](), Children: children})
}

type roleRepoFixture struct {
}

func (r *roleRepoFixture) Insert(rr RoleRoot) error {
	return nil
}
func (r *roleRepoFixture) SelectAll() ([]RoleRef, error) {
	return []RoleRef{}, nil
}
func (r *roleRepoFixture) SelectById(id id.ADT[ID]) (RoleRoot, error) {
	return RoleRoot{}, nil
}
func (r *roleRepoFixture) SelectChildren(id id.ADT[ID]) ([]RoleRef, error) {
	return []RoleRef{}, nil
}
func (r *roleRepoFixture) SelectState(id id.ADT[ID]) (state.Root, error) {
	return nil, nil
}

type stateRepoFixture struct {
}

func (r *stateRepoFixture) Insert(root state.Root) error {
	return nil
}
func (r *stateRepoFixture) SelectAll() ([]state.Ref, error) {
	return []state.Ref{}, nil
}
func (r *stateRepoFixture) SelectById(sid id.ADT[state.ID]) (state.Root, error) {
	return nil, nil
}
func (r *stateRepoFixture) SelectNext(sid id.ADT[state.ID]) (state.Ref, error) {
	return nil, nil
}

type kinshipRepoFixture struct {
}

func (r *kinshipRepoFixture) Insert(kr KinshipRoot) error {
	return nil
}
