package ph

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type Kind string

const (
	ID  = Kind("id")
	Sym = Kind("sym")
)

var kindRequired = []validation.Rule{
	validation.Required,
	validation.In(ID, Sym),
}

type Msg struct {
	K   Kind   `json:"kind"`
	ID  string `json:"id,omitempty"`
	Sym string `json:"sym,omitempty"`
}

func (dto Msg) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.K, kindRequired...),
		validation.Field(&dto.ID, id.RequiredWhen(dto.K == ID)...),
		validation.Field(&dto.Sym, sym.ReqiredWhen(dto.K == Sym)...),
	)
}

func MsgFromPH(ph ADT) Msg {
	switch val := ph.(type) {
	case id.ADT:
		return Msg{K: ID, ID: val.String()}
	case sym.ADT:
		return Msg{K: Sym, Sym: sym.ConvertToString(val)}
	default:
		panic(errUnexpectedType(ph))
	}
}

func MsgToPH(dto Msg) (ADT, error) {
	switch dto.K {
	case ID:
		return id.ConvertFromString(dto.ID)
	case Sym:
		return sym.CovertFromString(dto.Sym), nil
	default:
		panic(fmt.Errorf("unexpected placeholder kind: %v", dto.K))
	}
}
