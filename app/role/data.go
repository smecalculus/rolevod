package role

import (
	"database/sql"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type RoleRefData struct {
	ID   string         `db:"id"`
	Name string         `db:"name"`
	St   sql.NullString `db:"state"`
}

type roleRootData struct {
	ID       string         `db:"id"`
	Name     string         `db:"name"`
	St       sql.NullString `db:"state"`
	Children []RoleRefData  `db:"-"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
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
	id, err := id.StringTo(dto.ID)
	if err != nil {
		return RoleRoot{}, nil
	}
	state, err := state.JsonToRef(dto.St)
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
		St:       state,
		Children: children,
	}, nil
}

func dataFromRoleRoot(root RoleRoot) (roleRootData, error) {
	stateJson, err := state.JsonFromRef(root.St)
	if err != nil {
		return roleRootData{}, err
	}
	childrenDTOs, err := DataFromRoleRefs(root.Children)
	if err != nil {
		return roleRootData{}, err
	}
	return roleRootData{
		ID:       root.ID.String(),
		Name:     root.Name,
		St:       stateJson,
		Children: childrenDTOs,
	}, nil
}

type kinshipRootData struct {
	Parent   RoleRefData
	Children []RoleRefData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/internal/state:Json.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) (kinshipRootData, error)
)
