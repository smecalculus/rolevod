package cont

import (
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

type ID interface{}

type Ref interface {
}

type Root interface {
}

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectById(id.ADT[ID]) (Root, error)
	SelectByCh(id.ADT[chnl.ID]) (Root, error)
}
