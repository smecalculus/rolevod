package role

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/state"
)

type RefView struct {
	ID   string `form:"id" json:"id" param:"id"`
	Rev  int64  `form:"rev" json:"rev"`
	Name string `form:"name" json:"name"`
}

func (dto RefView) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
		validation.Field(&dto.Name, sym.Required...),
	)
}

type SpecView struct {
	NS   string `form:"ns"`
	Name string `form:"name"`
}

func (dto SpecView) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.NS, sym.Required...),
		validation.Field(&dto.Name, sym.Required...),
	)
}

type RootView struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	State state.SpecMsg `json:"state"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	ViewFromRef  func(Ref) RefView
	ViewToRef    func(RefView) (Ref, error)
	ViewFromRefs func([]Ref) []RefView
	ViewToRefs   func([]RefView) ([]Ref, error)
	ViewFromRoot func(Root) RootView
)
