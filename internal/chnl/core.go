package chnl

import (
	"errors"
	"fmt"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type Name = string
type ID = id.ADT

// aka ChanTp
type Spec struct {
	Name Name
	StID state.ID
	St   state.Ref
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
	// State
	St state.Ref
}

type Repo interface {
	Insert(Root) error
	InsertCtx([]Root) ([]Root, error)
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectCfg([]ID) (map[ID]Root, error)
	SelectMany([]ID) ([]Root, error)
}

func CollectStIDs(roots []Root) []state.ID {
	// stIDs := make([]state.ID, len(roots))
	var stIDs []state.ID
	for _, r := range roots {
		stIDs = append(stIDs, r.StID)
	}
	return stIDs
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Ident
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:Ident
var (
	ConvertRootToRef func(Root) Ref
)

var (
	ErrAlreadyConsumed = errors.New("channel already consumed")
)

func ErrDoesNotExist(rid ID) error {
	return fmt.Errorf("channel doesn't exist: %v", rid)
}

func ErrAlreadyClosed(rid ID) error {
	return fmt.Errorf("channel already closed: %v", rid)
}

func ErrNotAChnl(ph core.Placeholder) error {
	return fmt.Errorf("not a channel: %v", ph)
}
