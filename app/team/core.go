package team

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rev"
)

type ID = id.ADT
type Rev = rev.ADT
type Title = string

type Spec struct {
	Title string
	SupID id.ADT
}

type Ref struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
}

type Snap struct {
	ID    id.ADT
	Rev   rev.ADT
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
	RetreiveAll() ([]Ref, error)
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

type service struct {
	teams repo
	log   *slog.Logger
}

func newService(teams repo, l *slog.Logger) *service {
	name := slog.String("name", "teamService")
	return &service{teams, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	root := Root{
		ID:    id.New(),
		Rev:   rev.Initial(),
		Title: spec.Title,
		SupID: spec.SupID,
	}
	err := s.teams.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) Retrieve(rid id.ADT) (Snap, error) {
	snap, err := s.teams.SelectByID(rid)
	if err != nil {
		return Snap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveAll() ([]Ref, error) {
	return s.teams.SelectAll()
}

// Port
type repo interface {
	Insert(Root) error
	SelectByID(ID) (Snap, error)
	SelectAll() ([]Ref, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)
