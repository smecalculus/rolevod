package env

import (
	"log/slog"

	"smecalculus/rolevod/app/decl"
	"smecalculus/rolevod/lib/core"
)

type Spec struct {
	Name string
}

type Env core.Entity

type Root struct {
	ID    core.ID[Env]
	Name  string
	Decls []decl.TpDef
}

// port
type Api interface {
	Create(Spec) (Root, error)
	Retrieve(core.ID[Env]) (Root, error)
	RetreiveAll() ([]Root, error)
}

// core
type service struct {
	repo repo
	log  *slog.Logger
}

func newService(r repo, l *slog.Logger) *service {
	name := slog.String("name", "env.service")
	return &service{r, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	root := Root{
		ID:   core.New[Env](),
		Name: spec.Name,
	}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) Retrieve(id core.ID[Env]) (Root, error) {
	root, err := s.repo.SelectById(id)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) RetreiveAll() ([]Root, error) {
	roots, err := s.repo.SelectAll()
	if err != nil {
		return roots, err
	}
	return roots, nil
}

// port
type repo interface {
	Insert(Root) error
	SelectById(core.ID[Env]) (Root, error)
	SelectAll() ([]Root, error)
}

func toCore(id string) (core.ID[Env], error) {
	return core.FromString[Env](id)
}

func toEdge(id core.ID[Env]) string {
	return core.ToString(id)
}
