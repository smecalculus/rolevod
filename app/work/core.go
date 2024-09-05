package work

import (
	"log/slog"

	"smecalculus/rolevod/lib/core"
)

// Aggregate Root
type Work interface {
	work()
}

type WorkSpec struct {
	Name string
}

type WorkTeaser struct {
	ID   core.ID[Work]
	Name string
}

type WorkRoot struct {
	ID       core.ID[Work]
	Name     string
	Children []WorkTeaser
}

func (WorkRoot) work() {}

// Port
type WorkApi interface {
	Create(WorkSpec) (WorkRoot, error)
	Retrieve(core.ID[Work]) (WorkRoot, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]WorkTeaser, error)
}

type workService struct {
	workRepo    workRepo
	kinshipRepo kinshipRepo
	log         *slog.Logger
}

func newWorkService(sr workRepo, kr kinshipRepo, l *slog.Logger) *workService {
	name := slog.String("name", "workService")
	return &workService{sr, kr, l.With(name)}
}

func (s *workService) Create(spec WorkSpec) (WorkRoot, error) {
	root := WorkRoot{
		ID:   core.New[Work](),
		Name: spec.Name,
	}
	err := s.workRepo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *workService) Retrieve(id core.ID[Work]) (WorkRoot, error) {
	root, err := s.workRepo.SelectById(id)
	if err != nil {
		return WorkRoot{}, err
	}
	root.Children, err = s.workRepo.SelectChildren(id)
	if err != nil {
		return WorkRoot{}, err
	}
	return root, nil
}

func (s *workService) Establish(spec KinshipSpec) error {
	var children []WorkTeaser
	for _, id := range spec.ChildrenIDs {
		children = append(children, WorkTeaser{ID: id})
	}
	root := KinshipRoot{
		Parent:   WorkTeaser{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinshipRepo.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", root))
	return nil
}

func (s *workService) RetreiveAll() ([]WorkTeaser, error) {
	return s.workRepo.SelectAll()
}

// Port
type workRepo interface {
	Insert(WorkRoot) error
	SelectById(core.ID[Work]) (WorkRoot, error)
	SelectChildren(core.ID[Work]) ([]WorkTeaser, error)
	SelectAll() ([]WorkTeaser, error)
}

type KinshipSpec struct {
	ParentID    core.ID[Work]
	ChildrenIDs []core.ID[Work]
}

type KinshipRoot struct {
	Parent   WorkTeaser
	Children []WorkTeaser
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToWorkTeaser func(WorkRoot) WorkTeaser
)

func toSame(id core.ID[Work]) core.ID[Work] {
	return id
}

func toCore(id string) (core.ID[Work], error) {
	return core.FromString[Work](id)
}

func toEdge(id core.ID[Work]) string {
	return id.String()
}
