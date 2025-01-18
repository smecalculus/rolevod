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

func (dto SpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Name, core.NameRequired...),
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

type DealRefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

func (dto DealRefMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
		validation.Field(&dto.Name, core.NameRequired...),
	)
}

type RootMsg struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Sigs     []sig.RefMsg `json:"sigs"`
	Children []DealRefMsg `json:"children"`
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

type PartSpecMsg struct {
	Deal  string   `json:"deal_id" param:"id"`
	Sig   string   `json:"sig_id"`
	Owner string   `json:"owner_id"`
	TEs   []string `json:"tes"`
}

func (dto PartSpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Deal, id.Required...),
		validation.Field(&dto.Sig, id.Required...),
		validation.Field(&dto.TEs, core.CtxOptional...),
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
	Deal string       `json:"did"`
	PID  string       `json:"pid"`
	Key  string       `json:"key"`
	Term step.TermMsg `json:"term"`
}

func (dto TranSpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Deal, id.Required...),
		validation.Field(&dto.PID, id.Required...),
		validation.Field(&dto.Key, ak.Required...),
		validation.Field(&dto.Term, validation.Required),
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
