package agent

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"
)

type ID interface{}

type AgentSpec struct {
	Name string
}

type AgentRef struct {
	ID   id.ADT[ID]
	Name string
}

// Aggregate Root
type AgentRoot struct {
	ID       id.ADT[ID]
	Name     string
	Children []AgentRef
}

// Port
type AgentApi interface {
	Create(AgentSpec) (AgentRoot, error)
	Retrieve(id.ADT[ID]) (AgentRoot, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]AgentRef, error)
}

type agentService struct {
	agentRepo   agentRepo
	kinshipRepo kinshipRepo
	log         *slog.Logger
}

func newAgentService(sr agentRepo, kr kinshipRepo, l *slog.Logger) *agentService {
	name := slog.String("name", "agentService")
	return &agentService{sr, kr, l.With(name)}
}

func (s *agentService) Create(spec AgentSpec) (AgentRoot, error) {
	root := AgentRoot{
		ID:   id.New[ID](),
		Name: spec.Name,
	}
	err := s.agentRepo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *agentService) Retrieve(id id.ADT[ID]) (AgentRoot, error) {
	root, err := s.agentRepo.SelectById(id)
	if err != nil {
		return AgentRoot{}, err
	}
	root.Children, err = s.agentRepo.SelectChildren(id)
	if err != nil {
		return AgentRoot{}, err
	}
	return root, nil
}

func (s *agentService) Establish(spec KinshipSpec) error {
	var children []AgentRef
	for _, id := range spec.ChildrenIDs {
		children = append(children, AgentRef{ID: id})
	}
	root := KinshipRoot{
		Parent:   AgentRef{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinshipRepo.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", root))
	return nil
}

func (s *agentService) RetreiveAll() ([]AgentRef, error) {
	return s.agentRepo.SelectAll()
}

// Port
type agentRepo interface {
	Insert(AgentRoot) error
	SelectById(id.ADT[ID]) (AgentRoot, error)
	SelectChildren(id.ADT[ID]) ([]AgentRef, error)
	SelectAll() ([]AgentRef, error)
}

type KinshipSpec struct {
	ParentID    id.ADT[ID]
	ChildrenIDs []id.ADT[ID]
}

type KinshipRoot struct {
	Parent   AgentRef
	Children []AgentRef
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToAgentRef func(AgentRoot) AgentRef
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
