package chnl

import (
	"fmt"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"

	"smecalculus/rolevod/internal/state"
)

type Name = string
type ID = id.ADT

// aka ChanTp
type Spec struct {
	Name Name
	StID state.ID
}

// aka Z
type Ref struct {
	ID   ID
	Name Name
}

type Root struct {
	ID   ID
	Name Name
	// Preceding Channel ID
	PreID ID
	// Channel State ID
	StID state.ID
}

type Repo interface {
	Insert(Root) error
	InsertCtx([]Root) ([]Root, error)
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectByIDs([]ID) ([]Root, error)
	SelectCtx(ID, []ID) ([]Root, error)
	SelectCfg([]ID) (map[ID]Root, error)
	Transfer(from ID, to ID, pids []ID) error
}

func CollectStIDs(roots []Root) []state.ID {
	// stIDs := make([]state.ID, len(roots))
	var stIDs []state.ID
	for _, r := range roots {
		stIDs = append(stIDs, r.StID)
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
// goverter:extend smecalculus/rolevod/lib/id:Ident
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:Ident
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
