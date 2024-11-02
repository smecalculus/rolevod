package role

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"
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
	ID string `json:"id" param:"id"`
}

func (dto RefMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
	)
}

type RoleRefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

type RootMsg struct {
	ID       string        `json:"id" param:"id"`
	Name     string        `json:"name"`
	State    state.SpecMsg `json:"state"`
	Children []RoleRefMsg  `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgFromSpec func(Spec) SpecMsg
	MsgToSpec   func(SpecMsg) (Spec, error)
	MsgFromRef  func(Ref) RoleRefMsg
	MsgToRef    func(RoleRefMsg) (Ref, error)
	MsgFromRefs func([]Ref) []RoleRefMsg
	MsgToRefs   func([]RoleRefMsg) ([]Ref, error)
	MsgFromRoot func(Root) RootMsg
	// goverter:ignore StID
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoots func([]Root) []RootMsg
	MsgToRoots   func([]RootMsg) ([]Root, error)
)

type KinshipSpecMsg struct {
	ParentID string   `json:"parent_id" param:"id"`
	ChildIDs []string `json:"child_ids"`
}

type KinshipRootMsg struct {
	Parent   RoleRefMsg   `json:"parent"`
	Children []RoleRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
