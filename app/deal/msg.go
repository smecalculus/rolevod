package deal

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/app/sig"
	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/step"
)

type SpecMsg struct {
	Name string `json:"name"`
}

func (mto SpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Name, core.NameRequired...),
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

type DealRefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

func (mto DealRefMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.ID, id.Required...),
		validation.Field(&mto.Name, core.NameRequired...),
	)
}

type RootMsg struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Sigs     []sig.SigRefMsg `json:"sigs"`
	Children []DealRefMsg    `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/app/sig:Msg.*
var (
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromSpec  func(Spec) SpecMsg
	MsgToRef     func(DealRefMsg) (Ref, error)
	MsgFromRef   func(Ref) *DealRefMsg
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
)

type KinshipSpecMsg struct {
	ParentID string   `json:"parent_id" param:"id"`
	ChildIDs []string `json:"child_ids"`
}

func (mto KinshipSpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.ParentID, id.Required...),
		validation.Field(&mto.ChildIDs,
			validation.Required, validation.Length(0, 10),
			validation.Each(id.Required...)),
	)
}

type KinshipRootMsg struct {
	Parent   DealRefMsg   `json:"parent"`
	Children []DealRefMsg `json:"children"`
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

type PartSpecMsg struct {
	Deal  string   `json:"deal_id" param:"id"`
	Decl  string   `json:"sig_id"`
	Owner string   `json:"owner_id"`
	TEs   []string `json:"tes"`
}

func (mto PartSpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Deal, id.Required...),
		validation.Field(&mto.Decl, id.Required...),
		validation.Field(&mto.TEs, core.CtxOptional...),
	)
}

type PartRootMsg struct {
	ID     string        `json:"id"`
	DealID string        `json:"deal_id"`
	SigID  string        `json:"sig_id"`
	PID    string        `json:"pid"`
	TEs    []chnl.RefMsg `json:"tes"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/ak:Convert.*
// goverter:extend smecalculus/rolevod/internal/chnl:Msg.*
// goverter:extend smecalculus/rolevod/app/sig:Msg.*
var (
	MsgFromPartSpec func(PartSpec) PartSpecMsg
	MsgToPartSpec   func(PartSpecMsg) (PartSpec, error)
)

type TranSpecMsg struct {
	DID  string       `json:"did"`
	PID  string       `json:"pid"`
	Key  string       `json:"key"`
	Term step.TermMsg `json:"term"`
}

func (mto TranSpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.DID, id.Required...),
		validation.Field(&mto.PID, id.Required...),
		validation.Field(&mto.Key, ak.Required...),
		validation.Field(&mto.Term, validation.Required),
	)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/ak:Convert.*
// goverter:extend smecalculus/rolevod/internal/step:Msg.*
var (
	MsgFromTranSpec func(TranSpec) TranSpecMsg
	MsgToTranSpec   func(TranSpecMsg) (TranSpec, error)
)
