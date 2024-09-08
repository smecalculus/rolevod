package step

import (
	"errors"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

type ID interface{}

type Ref interface {
	get() id.ADT[ID]
}

type Label string

// aka Stype
type Root interface {
	Ref
}

type Fwd struct {
	ID   id.ADT[ID]
	From chnl.Ref
	To   chnl.Ref
}

func (r Fwd) get() id.ADT[ID] { return r.ID }

type Spawn struct {
	ID     id.ADT[ID]
	Name   string
	Via    chnl.Ref
	Values []chnl.Ref
	Cont   Root
}

func (r Spawn) get() id.ADT[ID] { return r.ID }

// aka ExpName
type ExpRef struct {
	ID     id.ADT[ID]
	Name   string
	Via    chnl.Ref
	Values []chnl.Ref
}

func (r ExpRef) get() id.ADT[ID] { return r.ID }

type Lab struct {
	ID   id.ADT[ID]
	Via  chnl.Ref
	Data Label
	// Cont Root
}

func (r Lab) get() id.ADT[ID] { return r.ID }

type Case struct {
	ID    id.ADT[ID]
	Via   chnl.Ref
	Conts map[Label]Root
}

func (r Case) get() id.ADT[ID] { return r.ID }

type Send struct {
	ID    id.ADT[ID]
	Via   chnl.Ref
	Value chnl.Ref
	// Cont  Root
}

func (r Send) get() id.ADT[ID] { return r.ID }

type Recv struct {
	ID    id.ADT[ID]
	Via   chnl.Ref
	Value chnl.Ref
	Cont  Root
}

func (r Recv) get() id.ADT[ID] { return r.ID }

type Close struct {
	ID  id.ADT[ID]
	Via chnl.Ref
}

func (r Close) get() id.ADT[ID] { return r.ID }

type Wait struct {
	ID   id.ADT[ID]
	Via  chnl.Ref
	Cont Root
}

func (r Wait) get() id.ADT[ID] { return r.ID }

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectById(id.ADT[ID]) (Root, error)
	SelectByChnl(chnl.ID) (Root, error)
}

var (
	ErrUnexpectedSt = errors.New("unexpected step type")
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToCoreIDs func([]string) ([]id.ADT[ID], error)
	ToEdgeIDs func([]id.ADT[ID]) []string
)

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
