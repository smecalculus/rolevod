package pool

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"
)

type SpecMsg struct {
	Title  string   `json:"title"`
	SupID  string   `json:"sup_id"`
	DepIDs []string `json:"dep_ids"`
}

func (dto SpecMsg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.SupID, id.Optional...),
	)
}

type IdentMsg struct {
	ID string `json:"id" param:"id"`
}

type RefMsg struct {
	ID    string `json:"id" param:"id"`
	Rev   int64  `json:"rev"`
	Title string `json:"title"`
}

type SnapMsg struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Subs  []RefMsg `json:"subs"`
}

type RootMsg struct {
	ID    string `json:"id"`
	Rev   int64  `json:"rev"`
	Title string `json:"title"`
	SupID string `json:"sup_id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromSpec  func(Spec) SpecMsg
	MsgToRoot    func(RootMsg) (Root, error)
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
	MsgToSnap    func(SnapMsg) (SubSnap, error)
	MsgFromSnap  func(SubSnap) SnapMsg
)
