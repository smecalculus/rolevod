package role

import (
	"errors"
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

var (
	ErrUnexpectedSt = errors.New("unexpected session type")
)

type ID interface{}

type RoleSpec struct {
	Name  string
	State state.Root
}

type RoleRef struct {
	ID   id.ADT[ID]
	Name string
}

// Aggregate Root
// aka TpDef
type RoleRoot struct {
	ID       id.ADT[ID]
	Name     string
	Children []RoleRef
	State    state.Root
}

// Port
type RoleApi interface {
	Create(RoleSpec) (RoleRoot, error)
	Retrieve(id.ADT[ID]) (RoleRoot, error)
	Update(RoleRoot) error
	Establish(KinshipSpec) error
	RetreiveAll() ([]RoleRef, error)
}

type roleService struct {
	roleRepo    roleRepo
	stateRepo   state.Repo
	kinshipRepo kinshipRepo
	log         *slog.Logger
}

func newRoleService(rr roleRepo, sr state.Repo, kr kinshipRepo, l *slog.Logger) *roleService {
	name := slog.String("name", "roleService")
	return &roleService{rr, sr, kr, l.With(name)}
}

func (s *roleService) Create(spec RoleSpec) (RoleRoot, error) {
	root := RoleRoot{
		ID:   id.New[ID](),
		Name: spec.Name,
	}
	if spec.State != nil {
		err := s.stateRepo.Insert(elab(spec.State))
		if err != nil {
			s.log.Error("creation failed", slog.Any("reason", err), slog.Any("state", spec.State))
			return root, err
		}
	}
	err := s.roleRepo.Insert(root)
	if err != nil {
		s.log.Error("creation failed", slog.Any("reason", err), slog.Any("spec", spec))
		return root, err
	}
	s.log.Debug("creation succeed", slog.Any("root", root))
	return root, nil
}

func (s *roleService) Update(root RoleRoot) error {
	return s.roleRepo.Insert(root)
}

func (s *roleService) Retrieve(rid id.ADT[ID]) (RoleRoot, error) {
	root, err := s.roleRepo.SelectById(rid)
	if err != nil {
		return RoleRoot{}, err
	}
	root.Children, err = s.roleRepo.SelectChildren(rid)
	if err != nil {
		return RoleRoot{}, err
	}
	// TODO: связка ролей и состояний
	// 1. переиспользуем идентификатор (/)
	// 2. складываем в отдельное поле
	// 3. строим отдельное отношение
	// - дополнительное обращение к БД
	stateID, err := id.String[state.ID](rid.String())
	if err != nil {
		return RoleRoot{}, err
	}
	root.State, err = s.stateRepo.SelectById(stateID)
	if err != nil {
		return RoleRoot{}, err
	}
	return root, nil
}

func (s *roleService) RetreiveAll() ([]RoleRef, error) {
	return s.roleRepo.SelectAll()
}

func (s *roleService) Establish(spec KinshipSpec) error {
	var children []RoleRef
	for _, id := range spec.Children {
		children = append(children, RoleRef{ID: id})
	}
	kr := KinshipRoot{
		Parent:   RoleRef{ID: spec.Parent},
		Children: children,
	}
	err := s.kinshipRepo.Insert(kr)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", kr))
	return nil
}

func elab(root state.Root) state.Root {
	switch st := root.(type) {
	case nil:
		return nil
	case *state.One:
		return &state.One{ID: id.New[state.ID]()}
	case *state.TpRef:
		return &state.TpRef{ID: id.New[state.ID](), Name: st.Name}
	default:
		panic(ErrUnexpectedSt)
	}
}

// Port
type roleRepo interface {
	Insert(RoleRoot) error
	SelectAll() ([]RoleRef, error)
	SelectById(id.ADT[ID]) (RoleRoot, error)
	SelectChildren(id.ADT[ID]) ([]RoleRef, error)
}

type KinshipSpec struct {
	Parent   id.ADT[ID]
	Children []id.ADT[ID]
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
var (
	ToRoleRef func(RoleRoot) RoleRef
	ToCoreIDs func([]string) ([]id.ADT[ID], error)
	ToEdgeIDs func([]id.ADT[ID]) []string
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
