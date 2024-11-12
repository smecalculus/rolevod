package alias

import (
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rev"
	"smecalculus/rolevod/lib/sym"
)

type Root struct {
	ID  id.ADT
	Rev rev.ADT
	Sym sym.ADT
}

type Repo interface {
	Insert(Root) error
}
