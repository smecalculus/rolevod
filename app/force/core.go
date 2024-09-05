package force

import (
	"log/slog"

	"smecalculus/rolevod/lib/core"
)

// Aggregate Root
type Force interface {
	force()
}

type ForceSpec struct {
	Name string
}

type ForceTeaser struct {
	ID   core.ID[Force]
	Name string
}

type ForceRoot struct {
	ID       core.ID[Force]
	Name     string
	Children []ForceTeaser
}

func (ForceRoot) force() {}

// Port
type ForceApi interface {
	Create(ForceSpec) (ForceRoot, error)
	Retrieve(core.ID[Force]) (ForceRoot, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]ForceTeaser, error)
}

type forceService struct {
	forceRepo   forceRepo
	kinshipRepo kinshipRepo
	log         *slog.Logger
}

func newForceService(sr forceRepo, kr kinshipRepo, l *slog.Logger) *forceService {
	name := slog.String("name", "forceService")
	return &forceService{sr, kr, l.With(name)}
}

func (s *forceService) Create(spec ForceSpec) (ForceRoot, error) {
	root := ForceRoot{
		ID:   core.New[Force](),
		Name: spec.Name,
	}
	err := s.forceRepo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *forceService) Retrieve(id core.ID[Force]) (ForceRoot, error) {
	root, err := s.forceRepo.SelectById(id)
	if err != nil {
		return ForceRoot{}, err
	}
	root.Children, err = s.forceRepo.SelectChildren(id)
	if err != nil {
		return ForceRoot{}, err
	}
	return root, nil
}

func (s *forceService) Establish(spec KinshipSpec) error {
	var children []ForceTeaser
	for _, id := range spec.ChildrenIDs {
		children = append(children, ForceTeaser{ID: id})
	}
	root := KinshipRoot{
		Parent:   ForceTeaser{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinshipRepo.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", root))
	return nil
}

func (s *forceService) RetreiveAll() ([]ForceTeaser, error) {
	return s.forceRepo.SelectAll()
}

// Port
type forceRepo interface {
	Insert(ForceRoot) error
	SelectById(core.ID[Force]) (ForceRoot, error)
	SelectChildren(core.ID[Force]) ([]ForceTeaser, error)
	SelectAll() ([]ForceTeaser, error)
}

type KinshipSpec struct {
	ParentID    core.ID[Force]
	ChildrenIDs []core.ID[Force]
}

type KinshipRoot struct {
	Parent   ForceTeaser
	Children []ForceTeaser
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToForceTeaser func(ForceRoot) ForceTeaser
)

func toSame(id core.ID[Force]) core.ID[Force] {
	return id
}

func toCore(id string) (core.ID[Force], error) {
	return core.FromString[Force](id)
}

func toEdge(id core.ID[Force]) string {
	return id.String()
}
