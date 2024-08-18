package env

import (
	"log/slog"

	"smecalculus/rolevod/app/dcl"
	"smecalculus/rolevod/lib/core"
)

// Aggregate Spec
type AS struct {
	Name string
}

// Aggregate Root
type AR struct {
	ID   core.ID[AR]
	Name string
	Tps  []dcl.TpRoot
	Exps []dcl.ExpRoot
}

// Port
type Api interface {
	Create(AS) (AR, error)
	Retrieve(core.ID[AR]) (AR, error)
	RetreiveAll() ([]AR, error)
}

type service struct {
	repo repo
	log  *slog.Logger
}

func newService(r repo, l *slog.Logger) *service {
	name := slog.String("name", "env.service")
	return &service{r, l.With(name)}
}

func (s *service) Create(spec AS) (AR, error) {
	root := AR{
		ID:   core.New[AR](),
		Name: spec.Name,
	}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) Retrieve(id core.ID[AR]) (AR, error) {
	root, err := s.repo.SelectById(id)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) RetreiveAll() ([]AR, error) {
	roots, err := s.repo.SelectAll()
	if err != nil {
		return roots, err
	}
	return roots, nil
}

// Port
type repo interface {
	Insert(AR) error
	SelectById(core.ID[AR]) (AR, error)
	SelectAll() ([]AR, error)
}

func toCore(id string) (core.ID[AR], error) {
	return core.FromString[AR](id)
}

func toEdge(id core.ID[AR]) string {
	return core.ToString(id)
}
