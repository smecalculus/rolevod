package sig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
)

type SpecMsg struct {
	FQN string         `json:"name"`
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

type RefMsg struct {
	ID string `json:"id" param:"id"`
}

func (dto RefMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
	)
}

type SigRefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name" query:"name"`
}

func (dto SigRefMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
	)
}

type RootMsg struct {
	ID   string         `json:"id"`
	Name string         `json:"name"`
	PE   chnl.SpecMsg   `json:"pe"`
	CEs  []chnl.SpecMsg `json:"ces"`
}

type SnapMsg struct {
	ID   string         `json:"id"`
	Name string         `json:"name"`
	PE   chnl.SpecMsg   `json:"pe"`
	CEs  []chnl.SpecMsg `json:"ces"`
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
	MsgToRef     func(SigRefMsg) (Ref, error)
	MsgFromRef   func(Ref) SigRefMsg
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
	MsgToSnap    func(SnapMsg) (Snap, error)
	MsgFromSnap  func(Snap) SnapMsg
)

type KinshipSpecMsg struct {
	ParentID string   `json:"parent_id" param:"id"`
	ChildIDs []string `json:"child_ids"`
}

type KinshipRootMsg struct {
	Parent   SigRefMsg   `json:"parent"`
	Children []SigRefMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
