package env

import (
	"log/slog"

	"smecalculus/rolevod/lib/core"
)

type Spec struct {
	Name string
}

type Id core.Kind

type Root struct {
	Id core.Id[Id]
}

// port
type Api interface {
	Create(Spec) (Root, error)
}

// core
type Service struct {
	repo Repo
	log  *slog.Logger
}

func NewService(r Repo, l *slog.Logger) *Service {
	name := slog.String("name", "env.Service")
	return &Service{r, l.With(name)}
}

func (s *Service) Create(spec Spec) (Root, error) {
	root := Root{}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

// port
type Repo interface {
	Insert(Root) error
}

func toCore(id string) (core.Id[Id], error) {
	return core.FromString[Id](id)
}

func toEdge(id core.Id[Id]) string {
	return core.ToString(id)
}
