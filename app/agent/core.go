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
	ID   id.ADT
	Name string
}

// Aggregate Root
type AgentRoot struct {
	ID       id.ADT
	Name     string
	Children []AgentRef
}

// Port
type AgentApi interface {
	Create(AgentSpec) (AgentRoot, error)
	Retrieve(id.ADT) (AgentRoot, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]AgentRef, error)
}

type agentService struct {
	agents   agentRepo
	kinships kinshipRepo
	log      *slog.Logger
}

func newAgentService(agents agentRepo, kinships kinshipRepo, l *slog.Logger) *agentService {
	name := slog.String("name", "agentService")
	return &agentService{agents, kinships, l.With(name)}
}

func (s *agentService) Create(spec AgentSpec) (AgentRoot, error) {
	root := AgentRoot{
		ID:   id.New(),
		Name: spec.Name,
	}
	err := s.agents.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *agentService) Retrieve(id id.ADT) (AgentRoot, error) {
	root, err := s.agents.SelectByID(id)
	if err != nil {
		return AgentRoot{}, err
	}
	root.Children, err = s.agents.SelectChildren(id)
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
	err := s.kinships.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeeded", slog.Any("kinship", root))
	return nil
}

func (s *agentService) RetreiveAll() ([]AgentRef, error) {
	return s.agents.SelectAll()
}

// Port
type agentRepo interface {
	Insert(AgentRoot) error
	SelectByID(id.ADT) (AgentRoot, error)
	SelectChildren(id.ADT) ([]AgentRef, error)
	SelectAll() ([]AgentRef, error)
}

type KinshipSpec struct {
	ParentID    id.ADT
	ChildrenIDs []id.ADT
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
// goverter:extend smecalculus/rolevod/lib/id:Ident
var (
	ToAgentRef func(AgentRoot) AgentRef
)
