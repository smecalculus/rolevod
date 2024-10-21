package role

import (
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type refData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type rootData struct {
	ID       string         `db:"id"`
	Name     string         `db:"name"`
	St       *state.RefData `db:"state"`
	Children []refData      `db:"-"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend data.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
var (
	DataToRef     func(refData) (Ref, error)
	DataFromRef   func(Ref) (refData, error)
	DataToRefs    func([]refData) ([]Ref, error)
	DataFromRefs  func([]Ref) ([]refData, error)
	DataToRoot    func(rootData) (Root, error)
	DataFromRoot  func(Root) (rootData, error)
	DataToRoots   func([]rootData) ([]Root, error)
	DataFromRoots func([]Root) ([]rootData, error)
)

func dataToRoot(dto rootData) (Root, error) {
	id, err := id.StringToID(dto.ID)
	if err != nil {
		return Root{}, nil
	}
	state, err := state.DataToRef(dto.St)
	if err != nil {
		return Root{}, nil
	}
	children, err := DataToRefs(dto.Children)
	if err != nil {
		return Root{}, nil
	}
	return Root{
		ID:       id,
		Name:     dto.Name,
		St:       state,
		Children: children,
	}, nil
}

func dataFromRoot(root Root) (rootData, error) {
	stateDTO := state.DataFromRef(root.St)
	childrenDTOs, err := DataFromRefs(root.Children)
	if err != nil {
		return rootData{}, err
	}
	return rootData{
		ID:       root.ID.String(),
		Name:     root.Name,
		St:       stateDTO,
		Children: childrenDTOs,
	}, nil
}

type kinshipRootData struct {
	Parent   refData
	Children []refData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
var (
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) (kinshipRootData, error)
)
