package alias

import (
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rev"
	"smecalculus/rolevod/lib/sym"
)

type Root struct {
	Sym sym.ADT
	ID  id.ADT
	Rev rev.ADT
}

type Repo interface {
	Insert(Root) error
}
