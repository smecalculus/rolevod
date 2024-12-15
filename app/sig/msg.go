package sig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
)

type SpecMsg struct {
	FQN string         `json:"fqn"`
	PE  chnl.SpecMsg   `json:"pe"`
	CEs []chnl.SpecMsg `json:"ces"`
}

func (dto SpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.FQN, sym.Required...),
		validation.Field(&dto.PE, validation.Required),
		validation.Field(&dto.CEs, core.CtxOptional...),
	)
}

type IdentMsg struct {
	ID string `json:"id" param:"id"`
}

func (dto IdentMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
	)
}

type RefMsg struct {
	ID    string `json:"id" param:"id"`
	Rev   int64  `json:"rev"`
	Title string `json:"title"`
}

func (dto RefMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
	)
}

type RootMsg struct {
	ID    string         `json:"id"`
	Rev   int64          `json:"rev"`
	Title string         `json:"title"`
	PE    chnl.SpecMsg   `json:"pe"`
	CEs   []chnl.SpecMsg `json:"ces"`
}

type SnapMsg struct {
	ID    string         `json:"id"`
	Rev   int64          `json:"rev"`
	Title string         `json:"title"`
	PE    chnl.SpecMsg   `json:"pe"`
	CEs   []chnl.SpecMsg `json:"ces"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/app/role:Msg.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToID      func(string) (id.ADT, error)
	MsgFromID    func(id.ADT) string
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromSpec  func(Spec) SpecMsg
	MsgToRef     func(RefMsg) (Ref, error)
	MsgFromRef   func(Ref) RefMsg
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
	MsgToSnap    func(SnapMsg) (Snap, error)
	MsgFromSnap  func(Snap) SnapMsg
)
