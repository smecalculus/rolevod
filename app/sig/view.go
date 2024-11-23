package sig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type SpecView struct {
	NS   string `form:"ns" json:"ns"`
	Name string `form:"name" json:"name"`
}

func (dto SpecView) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.NS, sym.Required...),
		validation.Field(&dto.Name, sym.Required...),
	)
}

type RefView struct {
	ID string `form:"id" json:"id" param:"id"`
	// Rev  int64  `form:"rev" json:"rev"`
	Name string `form:"name" json:"name"`
}

func (dto RefView) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.ID, id.Required...),
		// validation.Field(&dto.Rev, rev.Optional...),
		validation.Field(&dto.Name, sym.Required...),
	)
}

type RootView struct {
	ID string `json:"id"`
	// Rev  int64  `json:"rev"`
	Name string         `json:"name"`
	PE   chnl.SpecMsg   `json:"pe"`
	CEs  []chnl.SpecMsg `json:"ces"`
}

type SnapView struct {
	ID string `json:"id"`
	// Rev  int64  `json:"rev"`
	Name string `json:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/rev:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
// goverter:extend smecalculus/rolevod/internal/chnl:Msg.*
var (
	ViewFromRef  func(Ref) RefView
	ViewToRef    func(RefView) (Ref, error)
	ViewFromRefs func([]Ref) []RefView
	ViewToRefs   func([]RefView) ([]Ref, error)
	ViewFromSnap func(Snap) SnapView
	ViewFromRoot func(Root) RootView
)
