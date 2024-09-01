package env

import (
	"log/slog"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/dcl"
)

// Aggregate Spec
type AS interface {
	as()
}

func (EnvSpec) as() {}

type EnvSpec struct {
	Name string
}

// Aggregate Root
type AR interface {
	ar()
}

func (EnvRoot) ar() {}

type EnvRoot struct {
	ID   core.ID[AR]
	Name string
	Tps  []dcl.TpTeaser
	Exps []dcl.ExpRoot
}

type TpIntro struct {
	EnvID core.ID[AR]
	TpID  core.ID[dcl.AR]
}

// Port
type EnvApi interface {
	Create(EnvSpec) (EnvRoot, error)
	Retrieve(core.ID[AR]) (EnvRoot, error)
	RetreiveAll() ([]EnvRoot, error)
	Introduce(TpIntro) error
}

type service struct {
	envRepo envRepo
	tpRepo  tpRepo
	log     *slog.Logger
}

func newService(er envRepo, tr tpRepo, l *slog.Logger) *service {
	name := slog.String("name", "env.service")
	return &service{er, tr, l.With(name)}
}

func (s *service) Create(spec EnvSpec) (EnvRoot, error) {
	root := EnvRoot{
		ID:   core.New[AR](),
		Name: spec.Name,
	}
	err := s.envRepo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) Retrieve(id core.ID[AR]) (EnvRoot, error) {
	root, err := s.envRepo.SelectById(id)
	if err != nil {
		return root, err
	}
	root.Tps, err = s.tpRepo.SelectById(id)
	s.log.Error("tps", slog.Any("tps", root.Tps))
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) RetreiveAll() ([]EnvRoot, error) {
	return s.envRepo.SelectAll()
}

func (s *service) Introduce(intro TpIntro) error {
	return s.tpRepo.Insert(intro)
}

// Port
type envRepo interface {
	Insert(EnvRoot) error
	SelectById(core.ID[AR]) (EnvRoot, error)
	SelectAll() ([]EnvRoot, error)
}

// Port
type tpRepo interface {
	Insert(TpIntro) error
	SelectById(core.ID[AR]) ([]dcl.TpTeaser, error)
}

func toCore(id string) (core.ID[AR], error) {
	return core.FromString[AR](id)
}

func toEdge(id core.ID[AR]) string {
	return core.ToString(id)
}
