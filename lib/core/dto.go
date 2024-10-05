package core

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type phKind string

const (
	ID  = phKind("id")
	Sym = phKind("sym")
)

var phKindRequired = []validation.Rule{
	validation.Required,
	validation.In(ID, Sym),
}

type PlaceholderDTO struct {
	K   phKind `json:"k"`
	ID  string `json:"id,omitempty"`
	Sym string `json:"sym,omitempty"`
}

func (dto PlaceholderDTO) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.K, phKindRequired...),
		validation.Field(&dto.ID, id.RequiredWhen(dto.K == ID)...),
		validation.Field(&dto.Sym, sym.ReqiredWhen(dto.K == Sym)...),
	)
}

var NameOptional = []validation.Rule{
	validation.Length(1, 64),
	validation.Match(regexp.MustCompile("^[0-9A-Za-z_.-]*$")),
}

var NameRequired = append(NameOptional, validation.Required)

var CtxOptional = []validation.Rule{
	validation.Length(0, 10),
	validation.Each(validation.Required),
}

var CtxRequired = append(NameOptional, validation.Required)

func MsgFromPH(ph Placeholder) PlaceholderDTO {
	switch val := ph.(type) {
	case id.ADT:
		return PlaceholderDTO{K: ID, ID: val.String()}
	case sym.ADT:
		return PlaceholderDTO{K: Sym, Sym: sym.StringFromSym(val)}
	default:
		panic(errUnexpectedPHType(ph))
	}
}

func MsgToPH(mto PlaceholderDTO) (Placeholder, error) {
	switch mto.K {
	case ID:
		return id.StringToID(mto.ID)
	case Sym:
		return sym.StringToSym(mto.Sym), nil
	default:
		panic(errUnexpectedPHKind(mto.K))
	}
}

func errUnexpectedPHKind(k phKind) error {
	return fmt.Errorf("unexpected placeholder kind: %v", k)
}
