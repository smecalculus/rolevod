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

type RoleSpec struct {
	// Fully Qualified Name
	FQN FQN
	St  state.Spec
}

type RoleRef struct {
	ID   ID
	Name Name
}

// aka TpDef
type RoleRoot struct {
	ID       ID
	Name     Name
	St       state.Ref
	Children []RoleRef
}

type RoleApi interface {
	Create(RoleSpec) (RoleRoot, error)
	Retrieve(ID) (RoleRoot, error)
	RetreiveAll() ([]RoleRef, error)
	Update(RoleRoot) error
	Establish(KinshipSpec) error
}

type roleService struct {
	roles    roleRepo
	states   state.Repo
	aliases  alias.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newRoleService(
	roles roleRepo,
	states state.Repo,
	aliases alias.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *roleService {
	name := slog.String("name", "roleService")
	return &roleService{roles, states, aliases, kinships, l.With(name)}
}

func (s *roleService) Create(spec RoleSpec) (RoleRoot, error) {
	s.log.Debug("role creation started", slog.Any("spec", spec))
	st := state.ConvertSpecToRoot(spec.St)
	root := RoleRoot{
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

func (s *roleService) Update(root RoleRoot) error {
	return s.roles.Insert(root)
}

func (s *roleService) Retrieve(rid ID) (RoleRoot, error) {
	root, err := s.roles.SelectByID(rid)
	if err != nil {
		return RoleRoot{}, err
	}
	root.Children, err = s.roles.SelectChildren(rid)
	if err != nil {
		return RoleRoot{}, err
	}
	return root, nil
}

func (s *roleService) RetreiveAll() ([]RoleRef, error) {
	return s.roles.SelectAll()
}

func (s *roleService) Establish(spec KinshipSpec) error {
	var children []RoleRef
	for _, id := range spec.ChildIDs {
		children = append(children, RoleRef{ID: id})
	}
	root := KinshipRoot{
		Parent:   RoleRef{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinships.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("kinship establishment succeeded", slog.Any("root", root))
	return nil
}

type roleRepo interface {
	Insert(RoleRoot) error
	SelectAll() ([]RoleRef, error)
	SelectByID(ID) (RoleRoot, error)
	SelectChildren(ID) ([]RoleRef, error)
}

type KinshipSpec struct {
	ParentID ID
	ChildIDs []ID
}

type KinshipRoot struct {
	Parent   RoleRef
	Children []RoleRef
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Ident
// goverter:extend smecalculus/rolevod/internal/state:Convert.*
var (
	ConverRootToRef func(RoleRoot) RoleRef
)
