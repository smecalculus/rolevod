package alias

import (
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type Root struct {
	Sym sym.ADT
	ID  id.ADT
}

type Repo interface {
	Insert(Root) error
}
