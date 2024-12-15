package role

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rev"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/state"
)

type SpecMsg struct {
	FQN   string        `json:"fqn"`
	State state.SpecMsg `json:"state"`
}

func (dto SpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.FQN, sym.Required...),
		validation.Field(&dto.State, validation.Required),
	)
}

type RefMsg struct {
	ID    string `json:"id" param:"id"`
	Rev   int64  `json:"rev" query:"rev"`
	Title string `json:"title"`
}

func (dto RefMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
		validation.Field(&dto.Rev, rev.Optional...),
	)
}

type SnapMsg struct {
	ID    string        `json:"id" param:"id"`
	Rev   int64         `json:"rev" query:"rev"`
	Title string        `json:"title"`
	FQN   string        `json:"fqn"`
	State state.SpecMsg `json:"state"`
}

func (dto SnapMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
		validation.Field(&dto.Rev, rev.Optional...),
		validation.Field(&dto.State, validation.Required),
	)
}

type RootMsg struct {
	ID      string        `json:"id" param:"id"`
	Rev     int64         `json:"rev"`
	Title   string        `json:"title"`
	StateID string        `json:"state_id"`
	State   state.SpecMsg `json:"state"`
	Parts   []RefMsg      `json:"parts"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/rev:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgFromSpec func(Spec) SpecMsg
	MsgToSpec   func(SpecMsg) (Spec, error)
	MsgFromRef  func(Ref) RefMsg
	MsgToRef    func(RefMsg) (Ref, error)
	MsgFromRefs func([]Ref) []RefMsg
	MsgToRefs   func([]RefMsg) ([]Ref, error)
	// MsgFromRoot func(Root) RootMsg
	// goverter:ignore WholeID
	MsgToRoot func(RootMsg) (Root, error)
	// MsgFromRoots func([]Root) []RootMsg
	MsgToRoots   func([]RootMsg) ([]Root, error)
	MsgFromSnap  func(Snap) SnapMsg
	MsgToSnap    func(SnapMsg) (Snap, error)
	MsgFromSnaps func([]Snap) []SnapMsg
	MsgToSnaps   func([]SnapMsg) ([]Snap, error)
)
