package role

import (
	"log/slog"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/alias"
	"smecalculus/rolevod/internal/state"
)

func TestKinshipEstalish(t *testing.T) {
	newService(&roleRepoStub{}, &stateRepoStub{}, &aliasRepoStub{}, slog.Default())
}

type roleRepoStub struct {
}

func (r *roleRepoStub) Insert(rr Root) error {
	return nil
}
func (r *roleRepoStub) SelectAll() ([]Ref, error) {
	return []Ref{}, nil
}
func (r *roleRepoStub) SelectByID(id id.ADT) (Root, error) {
	return Root{}, nil
}
func (r *roleRepoStub) SelectByIDs(ids []id.ADT) ([]Root, error) {
	return []Root{}, nil
}
func (r *roleRepoStub) SelectEnv([]ID) (map[ID]Root, error) {
	return nil, nil
}
func (r *roleRepoStub) SelectParts(id id.ADT) ([]Ref, error) {
	return []Ref{}, nil
}
func (r *roleRepoStub) SelectState(id id.ADT) (state.Root, error) {
	return nil, nil
}

type stateRepoStub struct {
}

func (r *stateRepoStub) Insert(root state.Root) error {
	return nil
}
func (r *stateRepoStub) SelectAll() ([]state.Ref, error) {
	return []state.Ref{}, nil
}
func (r *stateRepoStub) SelectByID(sid id.ADT) (state.Root, error) {
	return nil, nil
}
func (r *stateRepoStub) SelectEnv(ids []id.ADT) (map[state.ID]state.Root, error) {
	return nil, nil
}
func (r *stateRepoStub) SelectByIDs(ids []id.ADT) ([]state.Root, error) {
	return nil, nil
}

type aliasRepoStub struct {
}

func (r *aliasRepoStub) Insert(ar alias.Root) error {
	return nil
}
