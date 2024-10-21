package role

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/state"
)

type SpecMsg struct {
	FQN string        `json:"fqn"`
	St  state.SpecMsg `json:"state"`
}

func (mto SpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.FQN, sym.Required...),
		validation.Field(&mto.St, validation.Required),
	)
}

type RefMsg struct {
	ID string `json:"id" param:"id"`
}

func (mto RefMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.ID, id.Required...),
	)
}

type RoleRefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

type RootMsg struct {
	ID       string       `json:"id" param:"id"`
	Name     string       `json:"name"`
	St       state.RefMsg `json:"state"`
	Children []RoleRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgFromSpec  func(Spec) SpecMsg
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromRef   func(Ref) RoleRefMsg
	MsgToRef     func(RoleRefMsg) (Ref, error)
	MsgFromRefs  func([]Ref) []RoleRefMsg
	MsgToRefs    func([]RoleRefMsg) ([]Ref, error)
	MsgFromRoot  func(Root) RootMsg
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
