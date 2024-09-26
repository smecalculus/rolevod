package role

import (
	"smecalculus/rolevod/internal/state"
)

type RoleSpecMsg struct {
	Name string         `json:"name"`
	St   *state.SpecMsg `json:"state"`
}

type RefMsg struct {
	ID string `json:"id" param:"id"`
}

type RoleRefMsg struct {
	ID   string        `json:"id" param:"id"`
	Name string        `json:"name"`
	St   *state.RefMsg `json:"state"`
}

type RoleRootMsg struct {
	ID       string        `json:"id" param:"id"`
	Name     string        `json:"name"`
	St       *state.RefMsg `json:"state"`
	Children []RoleRefMsg  `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgFromRoleSpec  func(RoleSpec) RoleSpecMsg
	MsgToRoleSpec    func(RoleSpecMsg) (RoleSpec, error)
	MsgFromRoleRef   func(RoleRef) RoleRefMsg
	MsgToRoleRef     func(RoleRefMsg) (RoleRef, error)
	MsgFromRoleRefs  func([]RoleRef) []RoleRefMsg
	MsgToRoleRefs    func([]RoleRefMsg) ([]RoleRef, error)
	MsgFromRoleRoot  func(RoleRoot) RoleRootMsg
	MsgToRoleRoot    func(RoleRootMsg) (RoleRoot, error)
	MsgFromRoleRoots func([]RoleRoot) []RoleRootMsg
	MsgToRoleRoots   func([]RoleRootMsg) ([]RoleRoot, error)
)

type KinshipSpecMsg struct {
	ParentID    string   `json:"parent_id" param:"id"`
	ChildrenIDs []string `json:"children_ids"`
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
