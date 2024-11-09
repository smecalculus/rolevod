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

func (mto Msg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.K, kindRequired...),
		validation.Field(&mto.ID, id.RequiredWhen(mto.K == ID)...),
		validation.Field(&mto.Sym, sym.ReqiredWhen(mto.K == Sym)...),
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

func MsgToPH(mto Msg) (ADT, error) {
	switch mto.K {
	case ID:
		return id.ConvertFromString(mto.ID)
	case Sym:
		return sym.CovertFromString(mto.Sym), nil
	default:
		panic(fmt.Errorf("unexpected placeholder kind: %v", mto.K))
	}
}
