package role

import (
	"smecalculus/rolevod/internal/state"
)

type RoleSpecMsg struct {
	Name  string         `json:"name"`
	State *state.RootMsg `json:"state,omitempty"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type RoleRootMsg struct {
	ID       string         `param:"id" json:"id"`
	Name     string         `json:"name"`
	Children []RoleRefMsg   `json:"children"`
	State    *state.RootMsg `json:"state,omitempty"`
}

type RoleRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToRoleSpec    func(RoleSpecMsg) (RoleSpec, error)
	MsgFromRoleSpec  func(RoleSpec) RoleSpecMsg
	MsgFromRoleRoot  func(RoleRoot) RoleRootMsg
	MsgToRoleRoot    func(RoleRootMsg) (RoleRoot, error)
	MsgFromRoleRoots func([]RoleRoot) []RoleRootMsg
	MsgToRoleRoots   func([]RoleRootMsg) ([]RoleRoot, error)
	MsgFromRoleRef   func(RoleRef) RoleRefMsg
	MsgToRoleRef     func(RoleRefMsg) (RoleRef, error)
	MsgFromRoleRefs  func([]RoleRef) []RoleRefMsg
	MsgToRoleRefs    func([]RoleRefMsg) ([]RoleRef, error)
)

type KinshipSpecMsg struct {
	Parent   string   `param:"id" json:"parent"`
	Children []string `json:"children"`
}

type KinshipRootMsg struct {
	Parent   RoleRefMsg   `json:"parent"`
	Children []RoleRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
