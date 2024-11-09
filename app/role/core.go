package role

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rev"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/alias"
	"smecalculus/rolevod/internal/state"
)

// for external readability
type ID = id.ADT
type Rev = rev.ADT
type FQN = sym.ADT
type Name = string

type Spec struct {
	// Fully Qualified Name
	FQN   sym.ADT
	State state.Spec
}

type Ref struct {
	ID   id.ADT
	Name string
}

// aka TpDef
type Root struct {
	ID   id.ADT
	Rev  rev.ADT
	Name string
	// specification relationship
	StateID state.ID
	State   state.Spec // read only
	// composition relationship
	WholeID id.ADT
	Parts   []Ref // read only
}

type API interface {
	Define(sym.ADT) (Ref, error)
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (Root, error)
	RetreiveAll() ([]Ref, error)
	Update(Root) error
}

type service struct {
	roles   repo
	states  state.Repo
	aliases alias.Repo
	log     *slog.Logger
}

func newService(
	roles repo,
	states state.Repo,
	aliases alias.Repo,
	l *slog.Logger,
) *service {
	name := slog.String("name", "roleService")
	return &service{roles, states, aliases, l.With(name)}
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
		return ConvertRootToRef(root), err
	}
	err = s.roles.Insert(root)
	if err != nil {
		s.log.Error("role insertion failed",
			slog.Any("reason", err),
			slog.Any("fqn", fqn),
		)
		return ConvertRootToRef(root), err
	}
	return ConvertRootToRef(root), nil
}

func (s *service) Create(spec Spec) (Root, error) {
	s.log.Debug("role creation started", slog.Any("spec", spec))
	st := state.ConvertSpecToRoot(spec.State)
	root := Root{
		ID:      id.New(),
		Rev:     rev.New(),
		Name:    spec.FQN.Name(),
		StateID: st.Ident(),
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
	root.State = spec.State
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
	root.State, err = s.states.SelectByID(root.StateID)
	if err != nil {
		s.log.Error("state selection failed")
		return Root{}, err
	}
	return root, nil
}

func (s *service) RetreiveAll() ([]Ref, error) {
	return s.roles.SelectAll()
}

type repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT) (Root, error)
	SelectChildren(ID) ([]Ref, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)
