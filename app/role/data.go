package role

import (
	"smecalculus/rolevod/lib/id"
)

type RoleRefData struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	State string `db:"state"`
}

type roleRootData struct {
	ID       string        `db:"id"`
	Name     string        `db:"name"`
	Children []RoleRefData `db:"-"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend data.*
// goverter:extend smecalculus/rolevod/internal/state:Json.*
var (
	DataToRoleRef     func(RoleRefData) (RoleRef, error)
	DataFromRoleRef   func(RoleRef) (RoleRefData, error)
	DataToRoleRefs    func([]RoleRefData) ([]RoleRef, error)
	DataFromRoleRefs  func([]RoleRef) ([]RoleRefData, error)
	DataToRoleRoot    func(roleRootData) (RoleRoot, error)
	DataFromRoleRoot  func(RoleRoot) (roleRootData, error)
	DataToRoleRoots   func([]roleRootData) ([]RoleRoot, error)
	DataFromRoleRoots func([]RoleRoot) ([]roleRootData, error)
)

func dataToRoleRoot(dto roleRootData) (RoleRoot, error) {
	id, err := id.String[ID](dto.ID)
	if err != nil {
		return RoleRoot{}, nil
	}
	children, err := DataToRoleRefs(dto.Children)
	if err != nil {
		return RoleRoot{}, nil
	}
	return RoleRoot{
		ID:       id,
		Name:     dto.Name,
		Children: children,
	}, nil
}

func dataFromRoleRoot(root RoleRoot) (roleRootData, error) {
	dtos, err := DataFromRoleRefs(root.Children)
	if err != nil {
		return roleRootData{}, err
	}
	return roleRootData{
		ID:       root.ID.String(),
		Name:     root.Name,
		Children: dtos,
	}, nil
}

type kinshipRootData struct {
	Parent   RoleRefData
	Children []RoleRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/internal/state:Json.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) (kinshipRootData, error)
)
