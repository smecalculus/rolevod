package role

import (
	"errors"
	"log/slog"

	"smecalculus/rolevod/lib/core"
)

var (
	ErrUnexpectedSt = errors.New("unexpected session type")
)

type Tpname = string
type Expname = string

// Aggregate Root
type Role interface {
	role()
}

type RoleSpec struct {
	Name Tpname
	St   Stype
}

type RoleTeaser struct {
	ID   core.ID[Role]
	Name Tpname
}

func (RoleRoot) role() {}

// aka TpDef
type RoleRoot struct {
	ID       core.ID[Role]
	Name     Tpname
	Children []RoleTeaser
	St       Stype
}

type Label string
type Choices map[Label]Stype

type Chan struct {
	V string
}

type Stype interface {
	stype()
}

func (One) stype()    {}
func (TpRef) stype()  {}
func (Tensor) stype() {}
func (Lolli) stype()  {}
func (With) stype()   {}
func (Plus) stype()   {}
func (Up) stype()     {}
func (Down) stype()   {}

// External Choice
type With struct {
	ID  core.ID[Role]
	Chs Choices
}

// Internal Choice
type Plus struct {
	ID  core.ID[Role]
	Chs Choices
}

type Tensor struct {
	ID core.ID[Role]
	S  Stype
	T  Stype
}

type Lolli struct {
	ID core.ID[Role]
	S  Stype
	T  Stype
}

type One struct {
	ID core.ID[Role]
}

// aka TpName
type TpRef struct {
	ID   core.ID[Role]
	Name Tpname
}

type Up struct {
	ID core.ID[Role]
	A  Stype
}

type Down struct {
	ID core.ID[Role]
	A  Stype
}

type ChanTp struct {
	X  Chan
	Tp Stype
}

// Port
type RoleApi interface {
	Create(RoleSpec) (RoleRoot, error)
	Retrieve(core.ID[Role]) (RoleRoot, error)
	Update(RoleRoot) error
	Establish(KinshipSpec) error
	RetreiveAll() ([]RoleTeaser, error)
}

type roleService struct {
	roleRepo    roleRepo
	kinshipRepo kinshipRepo
	log         *slog.Logger
}

func newRoleService(rr roleRepo, kr kinshipRepo, l *slog.Logger) *roleService {
	name := slog.String("name", "roleService")
	return &roleService{rr, kr, l.With(name)}
}

func (s *roleService) Create(spec RoleSpec) (RoleRoot, error) {
	root := RoleRoot{
		ID:   core.New[Role](),
		Name: spec.Name,
		St:   elab(spec.St),
	}
	err := s.roleRepo.Insert(root)
	if err != nil {
		s.log.Error("creation failed", slog.Any("reason", err), slog.Any("spec", spec))
		return root, err
	}
	s.log.Debug("creation succeed", slog.Any("root", root))
	return root, nil
}

func (s *roleService) Update(rr RoleRoot) error {
	return s.roleRepo.Insert(rr)
}

func (s *roleService) Retrieve(id core.ID[Role]) (RoleRoot, error) {
	root, err := s.roleRepo.SelectById(id)
	if err != nil {
		return RoleRoot{}, err
	}
	root.Children, err = s.roleRepo.SelectChildren(id)
	if err != nil {
		return RoleRoot{}, err
	}
	return root, nil
}

func (s *roleService) RetreiveAll() ([]RoleTeaser, error) {
	return s.roleRepo.SelectAll()
}

func (s *roleService) Establish(ks KinshipSpec) error {
	var children []RoleTeaser
	for _, id := range ks.Children {
		children = append(children, RoleTeaser{ID: id})
	}
	kr := KinshipRoot{
		Parent:   RoleTeaser{ID: ks.Parent},
		Children: children,
	}
	err := s.kinshipRepo.Insert(kr)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", kr))
	return nil
}

func elab(stype Stype) Stype {
	switch st := stype.(type) {
	case nil:
		return nil
	case One:
		return One{ID: core.New[Role]()}
	case TpRef:
		return TpRef{ID: core.New[Role](), Name: st.Name}
	default:
		panic(ErrUnexpectedSt)
	}
}

// Port
type roleRepo interface {
	Insert(RoleRoot) error
	SelectById(core.ID[Role]) (RoleRoot, error)
	SelectChildren(core.ID[Role]) ([]RoleTeaser, error)
	SelectAll() ([]RoleTeaser, error)
}

type KinshipSpec struct {
	Parent   core.ID[Role]
	Children []core.ID[Role]
}

type KinshipRoot struct {
	Parent   RoleTeaser
	Children []RoleTeaser
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToRoleTeaser func(RoleRoot) RoleTeaser
	ToCoreIDs    func([]string) ([]core.ID[Role], error)
	ToEdgeIDs    func([]core.ID[Role]) []string
)

func toSame(id core.ID[Role]) core.ID[Role] {
	return id
}

func toCore(id string) (core.ID[Role], error) {
	return core.FromString[Role](id)
}

func toEdge(id core.ID[Role]) string {
	return id.String()
}
