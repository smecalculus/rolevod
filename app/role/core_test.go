package role

import (
	"log/slog"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/alias"
	"smecalculus/rolevod/internal/state"
)

func TestKinshipEstalish(t *testing.T) {
	s := newRoleService(&roleRepoFixture{}, &stateRepoFixture{}, &aliasRepoFixture{}, &kinshipRepoFixture{}, slog.Default())
	children := []id.ADT{id.New()}
	s.Establish(KinshipSpec{ParentID: id.New(), ChildIDs: children})
}

type roleRepoFixture struct {
}

func (r *roleRepoFixture) Insert(rr RoleRoot) error {
	return nil
}
func (r *roleRepoFixture) SelectAll() ([]RoleRef, error) {
	return []RoleRef{}, nil
}
func (r *roleRepoFixture) SelectByID(id id.ADT) (RoleRoot, error) {
	return RoleRoot{}, nil
}
func (r *roleRepoFixture) SelectChildren(id id.ADT) ([]RoleRef, error) {
	return []RoleRef{}, nil
}
func (r *roleRepoFixture) SelectState(id id.ADT) (state.Root, error) {
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
func (r *stateRepoFixture) SelectByID(sid id.ADT) (state.Root, error) {
	return nil, nil
}
func (r *stateRepoFixture) SelectEnv(ids []id.ADT) (map[state.ID]state.Root, error) {
	return nil, nil
}
func (r *stateRepoFixture) SelectByIDs(ids []id.ADT) ([]state.Root, error) {
	return nil, nil
}

type aliasRepoFixture struct {
}

func (r *aliasRepoFixture) Insert(ar alias.Root) error {
	return nil
}

type kinshipRepoFixture struct {
}

func (r *kinshipRepoFixture) Insert(kr KinshipRoot) error {
	return nil
}
