package decl

import (
	"log/slog"

	"smecalculus/rolevod/lib/core"
)

type Spec struct {
	Name string
}

type Decl core.Entity

type Root struct {
	ID   core.ID[Decl]
	Name string
}

// port
type Api interface {
	Create(Spec) (Root, error)
	Retrieve(core.ID[Decl]) (Root, error)
	RetreiveAll() ([]Root, error)
}

// core
type service struct {
	repo repo
	log  *slog.Logger
}

func newService(r repo, l *slog.Logger) *service {
	name := slog.String("name", "decl.service")
	return &service{r, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	root := Root{
		ID:   core.New[Decl](),
		Name: spec.Name,
	}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) Retrieve(id core.ID[Decl]) (Root, error) {
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
	SelectById(core.ID[Decl]) (Root, error)
	SelectAll() ([]Root, error)
}

func ToCore(id string) (core.ID[Decl], error) {
	return core.FromString[Decl](id)
}

func ToEdge(id core.ID[Decl]) string {
	return core.ToString(id)
}
