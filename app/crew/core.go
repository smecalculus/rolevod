package crew

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"
)

type ID = id.ADT

type Spec struct {
	Name string
}

type Ref struct {
	ID   ID
	Name string
}

type Root struct {
	ID       ID
	Name     string
	Children []Ref
}

// Port
type API interface {
	Create(Spec) (Root, error)
	Retrieve(ID) (Root, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]Ref, error)
}

type service struct {
	agents   repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newService(agents repo, kinships kinshipRepo, l *slog.Logger) *service {
	name := slog.String("name", "crewService")
	return &service{agents, kinships, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	root := Root{
		ID:   id.New(),
		Name: spec.Name,
	}
	err := s.agents.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) Retrieve(id ID) (Root, error) {
	root, err := s.agents.SelectByID(id)
	if err != nil {
		return Root{}, err
	}
	root.Children, err = s.agents.SelectChildren(id)
	if err != nil {
		return Root{}, err
	}
	return root, nil
}

func (s *service) Establish(spec KinshipSpec) error {
	var children []Ref
	for _, id := range spec.ChildIDs {
		children = append(children, Ref{ID: id})
	}
	root := KinshipRoot{
		Parent:   Ref{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinships.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeeded", slog.Any("kinship", root))
	return nil
}

func (s *service) RetreiveAll() ([]Ref, error) {
	return s.agents.SelectAll()
}

// Port
type repo interface {
	Insert(Root) error
	SelectByID(ID) (Root, error)
	SelectChildren(ID) ([]Ref, error)
	SelectAll() ([]Ref, error)
}

type KinshipSpec struct {
	ParentID ID
	ChildIDs []ID
}

type KinshipRoot struct {
	Parent   Ref
	Children []Ref
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ToAgentRef func(Root) Ref
)
