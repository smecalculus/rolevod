package ph

import (
	"fmt"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type kind int

const (
	nonph = iota
	idKind
	symKind
)

type Data struct {
	K   kind   `json:"k"`
	ID  string `json:"id,omitempty"`
	Sym string `json:"sym,omitempty"`
}

func DataFromPH(ph ADT) Data {
	switch val := ph.(type) {
	case id.ADT:
		return Data{K: idKind, ID: val.String()}
	case sym.ADT:
		return Data{K: symKind, Sym: sym.ConvertToString(val)}
	default:
		panic(errUnexpectedType(ph))
	}
}

func DataToPH(dto Data) (ADT, error) {
	switch dto.K {
	case idKind:
		return id.ConvertFromString(dto.ID)
	case symKind:
		return sym.CovertFromString(dto.Sym), nil
	default:
		panic(fmt.Errorf("unexpected placeholder kind: %v", dto.K))
	}
}
