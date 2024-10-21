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

func (mto SpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.FQN, sym.Required...),
		validation.Field(&mto.PE, validation.Required),
		validation.Field(&mto.CEs, core.CtxOptional...),
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

type SigRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

type RootMsg struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	PE       chnl.SpecMsg   `json:"pe"`
	CEs      []chnl.SpecMsg `json:"ces"`
	Children []SigRefMsg    `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
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
// goverter:extend smecalculus/rolevod/lib/id:String.*
var (
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
