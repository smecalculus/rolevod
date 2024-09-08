package role

import (
	"smecalculus/rolevod/lib/id"
)

type roleRootData struct {
	ID       string        `db:"id"`
	Name     string        `db:"name"`
	Children []roleRefData `db:"-"`
}

type roleRefData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend data.*
var (
	DataToRoleRef     func(roleRefData) (RoleRef, error)
	DataFromRoleRef   func(RoleRef) roleRefData
	DataToRoleRefs    func([]roleRefData) ([]RoleRef, error)
	DataFromRoleRefs  func([]RoleRef) []roleRefData
	DataToRoleRoot    func(roleRootData) (RoleRoot, error)
	DataFromRoleRoot  func(RoleRoot) roleRootData
	DataToRoleRoots   func([]roleRootData) ([]RoleRoot, error)
	DataFromRoleRoots func([]RoleRoot) []roleRootData
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

func dataFromRoleRoot(root RoleRoot) roleRootData {
	return roleRootData{
		ID:       root.ID.String(),
		Name:     root.Name,
		Children: DataFromRoleRefs(root.Children),
	}
}

type kinshipRootData struct {
	Parent   roleRefData
	Children []roleRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)
