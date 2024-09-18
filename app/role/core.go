package role

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type ID interface{}

type RoleSpec struct {
	Name string
	St   state.Spec
}

type RoleRef struct {
	ID   id.ADT[ID]
	Name string
	St   state.Ref
}

// aka TpDef
type RoleRoot struct {
	ID       id.ADT[ID]
	Name     string
	St       state.Ref
	Children []RoleRef
}

type RoleApi interface {
	Create(RoleSpec) (RoleRoot, error)
	Retrieve(id.ADT[ID]) (RoleRoot, error)
	RetreiveAll() ([]RoleRef, error)
	Update(RoleRoot) error
	Establish(KinshipSpec) error
}

type roleService struct {
	roles    roleRepo
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newRoleService(
	roles roleRepo,
	states state.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *roleService {
	name := slog.String("name", "roleService")
	return &roleService{roles, states, kinships, l.With(name)}
}

func (s *roleService) Create(spec RoleSpec) (RoleRoot, error) {
	s.log.Debug("role creation started", slog.Any("spec", spec))
	st := state.ConvertSpecToRoot(spec.St)
	root := RoleRoot{
		ID:   id.New[ID](),
		Name: spec.Name,
		// State: state.ConvertRootToRef(st),
		St: st,
	}
	err := s.roles.Insert(root)
	if err != nil {
		s.log.Error("role insertion failed",
			slog.Any("reason", err),
			slog.Any("spec", spec),
		)
		return root, err
	}
	if spec.St != nil {
		err := s.states.Insert(st)
		if err != nil {
			s.log.Error("state insertion failed",
				slog.Any("reason", err),
				slog.Any("root", st),
			)
			return root, err
		}
	}
	s.log.Debug("role creation succeeded", slog.Any("root", root))
	return root, nil
}

func (s *roleService) Update(root RoleRoot) error {
	return s.roles.Insert(root)
}

func (s *roleService) Retrieve(rid id.ADT[ID]) (RoleRoot, error) {
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
	for _, id := range spec.ChildrenIDs {
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
	SelectByID(id.ADT[ID]) (RoleRoot, error)
	SelectChildren(id.ADT[ID]) ([]RoleRef, error)
}

type KinshipSpec struct {
	ParentID    id.ADT[ID]
	ChildrenIDs []id.ADT[ID]
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
// goverter:extend to.*
// goverter:extend smecalculus/rolevod/internal/state:Convert.*
var (
	ConverRootToRef func(RoleRoot) RoleRef
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
