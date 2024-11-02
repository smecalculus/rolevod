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
	FQN   FQN
	State state.Spec
}

type Ref struct {
	ID   ID
	Name Name
}

// aka TpDef
type Root struct {
	ID   ID
	Name Name
	// Rev      int8
	StID     state.ID
	State    state.Root
	Children []Ref
}

type API interface {
	Define(FQN) (Ref, error)
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

func (s *service) Define(fqn FQN) (Ref, error) {
	s.log.Debug("role definition started", slog.Any("fqn", fqn))
	root := Root{
		ID:   id.New(),
		Name: fqn.Name(),
	}
	al := alias.Root{Sym: fqn, ID: root.ID}
	err := s.aliases.Insert(al)
	if err != nil {
		s.log.Error("alias insertion failed",
			slog.Any("reason", err),
			slog.Any("alias", al),
		)
		return ConverRootToRef(root), err
	}
	err = s.roles.Insert(root)
	if err != nil {
		s.log.Error("role insertion failed",
			slog.Any("reason", err),
			slog.Any("fqn", fqn),
		)
		return ConverRootToRef(root), err
	}
	return ConverRootToRef(root), nil
}

func (s *service) Create(spec Spec) (Root, error) {
	s.log.Debug("role creation started", slog.Any("spec", spec))
	st := state.ConvertSpecToRoot(spec.State)
	root := Root{
		ID:    id.New(),
		Name:  spec.FQN.Name(),
		StID:  st.Ident(),
		State: st,
	}
	err := s.roles.Insert(root)
	if err != nil {
		s.log.Error("role insertion failed",
			slog.Any("reason", err),
			slog.Any("root", root),
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
	alias := alias.Root{Sym: spec.FQN, ID: root.ID}
	err = s.aliases.Insert(alias)
	if err != nil {
		s.log.Error("alias insertion failed",
			slog.Any("reason", err),
			slog.Any("alias", alias),
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
		s.log.Error("root selection failed")
		return Root{}, err
	}
	root.State, err = s.states.SelectByID(root.StID)
	if err != nil {
		s.log.Error("state selection failed")
		return Root{}, err
	}
	// root.Children, err = s.roles.SelectChildren(rid)
	// if err != nil {
	// 	s.log.Error("children selection failed")
	// 	return Root{}, err
	// }
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
