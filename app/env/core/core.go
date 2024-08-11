package env

import (
	"log/slog"

	"smecalculus/rolevod/lib/core"
)

// message
type Spec struct {
	Id string `json:"id"`
}

type Id core.Kind

// domain
type Env struct {
	Id core.Id[Id]
}

// state
type Root struct {
	Id string
}

// port
type Api interface {
	Create(Spec) (Env, error)
}

// core
type Service struct {
	repo Repo
	log  *slog.Logger
}

func NewService(repo Repo, log *slog.Logger) *Service {
	name := slog.String("name", "env.Service")
	return &Service{repo, log.With(name)}
}

func (s *Service) Create(spec Spec) (Env, error) {
	root := Root{}
	err := s.repo.Insert(root)
	if err != nil {
		return Env{}, err
	}
	return Env{}, nil
}

// port
type Repo interface {
	Insert(Root) error
}
