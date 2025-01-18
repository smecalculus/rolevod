package pool

import (
	"log/slog"

	"smecalculus/rolevod/app/sig"
	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rev"
)

type ID = id.ADT
type Rev = rev.ADT
type Title = string

type Spec struct {
	Title  string
	SupID  id.ADT
	DepIDs []sig.ID
}

type Ref struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
}

type Snap struct {
	ID    id.ADT
	Title string
	Subs  []Ref
}

type Root struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
	SupID id.ADT
}

// Port
type API interface {
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (Snap, error)
	RetreiveRefs() ([]Ref, error)
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

type service struct {
	pools Repo
	log   *slog.Logger
}

func newService(pools Repo, l *slog.Logger) *service {
	name := slog.String("name", "poolService")
	return &service{pools, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	root := Root{
		ID:    id.New(),
		Rev:   rev.Initial(),
		Title: spec.Title,
		SupID: spec.SupID,
	}
	err := s.pools.Insert(root)
	if err != nil {
		return Root{}, err
	}
	return root, nil
}

func (s *service) Retrieve(rid id.ADT) (Snap, error) {
	snap, err := s.pools.SelectByID(rid)
	if err != nil {
		return Snap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveRefs() ([]Ref, error) {
	return s.pools.SelectAll()
}

// Port
type Repo interface {
	Insert(Root) error
	SelectByID(id.ADT) (Snap, error)
	SelectAll() ([]Ref, error)
	Transfer(giver id.ADT, taker id.ADT, pids []chnl.ID) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)
