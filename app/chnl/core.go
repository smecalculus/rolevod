package chnl

import (
	"fmt"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/rev"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/state"
)

// for external readability
type ID = id.ADT
type Key = string

// aka ChanTp
type Spec struct {
	Key  string
	Link sym.ADT
}

// aka Z
type Ref struct {
	ID  id.ADT
	Key string
}

type Root struct {
	ID      id.ADT
	Key     string
	Rev     rev.ADT
	PreID   *id.ADT
	StateID *state.ID
}

type Repo interface {
	Insert(Root) error
	InsertCtx([]Root) ([]Root, error)
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT) (Root, error)
	SelectByIDs([]id.ADT) ([]Root, error)
	SelectCtx(id.ADT, []id.ADT) ([]Root, error)
	SelectCfg([]id.ADT) (map[id.ADT]Root, error)
	Transfer(giver id.ADT, taker id.ADT, pids []id.ADT) error
}

func CollectCtx(roots []Root) []state.ID {
	var stIDs []state.ID
	for _, r := range roots {
		if r.StateID == nil {
			continue
		}
		stIDs = append(stIDs, *r.StateID)
	}
	return stIDs
}

func Subst(cur ID, pat ID, new ID) ID {
	if cur == pat {
		return new
	}
	return cur
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/ak:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)

func ErrDoesNotExist(want ID) error {
	return fmt.Errorf("channel doesn't exist: %v", want)
}

func ErrMissingInCfg(want ph.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func ErrMissingInCtx(want ph.ADT) error {
	return fmt.Errorf("channel missing in ctx: %v", want)
}

func ErrAlreadyClosed(got ID) error {
	return fmt.Errorf("channel already closed: %v", got)
}

func ErrNotAnID(got ph.ADT) error {
	return fmt.Errorf("not a channel id: %v", got)
}
