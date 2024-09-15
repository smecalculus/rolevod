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
		ID:   id.New[ID](),
		Name: spec.Name,
	}
	err := s.agents.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *agentService) Retrieve(id id.ADT[ID]) (AgentRoot, error) {
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
	SelectByID(id.ADT[ID]) (AgentRoot, error)
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
