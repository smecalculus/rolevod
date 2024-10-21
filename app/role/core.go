package role

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/alias"
	"smecalculus/rolevod/internal/state"
)

type ID = id.ADT
type FQN = sym.ADT
type Name = string

type Spec struct {
	// Fully Qualified Name
	FQN FQN
	St  state.Spec
}

type Ref struct {
	ID   ID
	Name Name
}

// aka TpDef
type Root struct {
	ID       ID
	Name     Name
	St       state.Ref
	Children []Ref
}

type API interface {
	Create(Spec) (Root, error)
	Retrieve(ID) (Root, error)
	RetreiveAll() ([]Ref, error)
	Update(Root) error
	Establish(KinshipSpec) error
}

type service struct {
	roles    repo
	states   state.Repo
	aliases  alias.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newService(
	roles repo,
	states state.Repo,
	aliases alias.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *service {
	name := slog.String("name", "roleService")
	return &service{roles, states, aliases, kinships, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	s.log.Debug("role creation started", slog.Any("spec", spec))
	st := state.ConvertSpecToRoot(spec.St)
	root := Root{
		ID:   id.New(),
		Name: spec.FQN.Name(),
		St:   st,
	}
	err := s.roles.Insert(root)
	if err != nil {
		s.log.Error("role insertion failed",
			slog.Any("reason", err),
			slog.Any("spec", spec),
		)
		return root, err
	}
	err = s.states.Insert(st)
	if err != nil {
		s.log.Error("state insertion failed",
			slog.Any("reason", err),
			slog.Any("state", st),
		)
		return root, err
	}
	al := alias.Root{Sym: spec.FQN, ID: root.ID}
	err = s.aliases.Insert(al)
	if err != nil {
		s.log.Error("alias insertion failed",
			slog.Any("reason", err),
			slog.Any("alias", al),
		)
		return root, err
	}
	s.log.Debug("role creation succeeded", slog.Any("root", root))
	return root, nil
}

func (s *service) Update(root Root) error {
	return s.roles.Insert(root)
}

func (s *service) Retrieve(rid ID) (Root, error) {
	root, err := s.roles.SelectByID(rid)
	if err != nil {
		return Root{}, err
	}
	root.Children, err = s.roles.SelectChildren(rid)
	if err != nil {
		return Root{}, err
	}
	return root, nil
}

func (s *service) RetreiveAll() ([]Ref, error) {
	return s.roles.SelectAll()
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
	s.log.Debug("kinship establishment succeeded", slog.Any("root", root))
	return nil
}

type repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectChildren(ID) ([]Ref, error)
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
// goverter:extend smecalculus/rolevod/lib/id:Ident
// goverter:extend smecalculus/rolevod/internal/state:Convert.*
var (
	ConverRootToRef func(Root) Ref
)
