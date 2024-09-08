package chnl

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"
)

type ID interface{}

type Spec struct {
	Name string
}

type Ref struct {
	ID   id.ADT[ID]
	Name string
}

// Aggregate Root
type Root struct {
	ID   id.ADT[ID]
	Name string
}

// Port
type Api interface {
	Create(Spec) (Root, error)
	Retrieve(id.ADT[ID]) (Root, error)
	RetreiveAll() ([]Ref, error)
}

type service struct {
	repo repo
	log  *slog.Logger
}

func newService(r repo, l *slog.Logger) *service {
	name := slog.String("name", "chnl.service")
	return &service{r, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	root := Root{
		ID:   id.New[ID](),
		Name: spec.Name,
	}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) Retrieve(id id.ADT[ID]) (Root, error) {
	return s.repo.SelectById(id)
}

func (s *service) RetreiveAll() ([]Ref, error) {
	return s.repo.SelectAll()
}

// Port
type repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectById(id.ADT[ID]) (Root, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToRef func(Root) Ref
)

func toSame(id id.ADT[ID]) id.ADT[ID] {
	return id
}

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
