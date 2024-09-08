package role

import (
	"errors"
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/internal/chnl"
)

var (
	ErrUnexpectedSt = errors.New("unexpected session type")
)

type Tpname = string
type Expname = string

// Aggregate Root
type ID interface{}

type RoleSpec struct {
	Name Tpname
	St   Stype
}

type RoleRef struct {
	ID   id.ADT[ID]
	Name Tpname
}

// Aggregate Root
// aka TpDef
type RoleRoot struct {
	ID       id.ADT[ID]
	Name     Tpname
	Children []RoleRef
	St       Stype
}

type Label string

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

// aka External Choice
type With struct {
	ID      id.ADT[ID]
	Choices map[Label]Stype
}

// aka Internal Choice
type Plus struct {
	ID      id.ADT[ID]
	Choices map[Label]Stype
}

type Tensor struct {
	ID id.ADT[ID]
	S  Stype
	T  Stype
}

type Lolli struct {
	ID id.ADT[ID]
	S  Stype
	T  Stype
}

type One struct {
	ID id.ADT[ID]
}

// aka TpName
type TpRef struct {
	ID   id.ADT[ID]
	Name Tpname
}

type Up struct {
	ID id.ADT[ID]
	A  Stype
}

type Down struct {
	ID id.ADT[ID]
	A  Stype
}

type ChanTp struct {
	X  chnl.Ref
	Tp Stype
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
	kinshipRepo kinshipRepo
	log         *slog.Logger
}

func newRoleService(rr roleRepo, kr kinshipRepo, l *slog.Logger) *roleService {
	name := slog.String("name", "roleService")
	return &roleService{rr, kr, l.With(name)}
}

func (s *roleService) Create(spec RoleSpec) (RoleRoot, error) {
	root := RoleRoot{
		ID:   id.New[ID](),
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

func (s *roleService) Retrieve(id id.ADT[ID]) (RoleRoot, error) {
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

func (s *roleService) RetreiveAll() ([]RoleRef, error) {
	return s.roleRepo.SelectAll()
}

func (s *roleService) Establish(ks KinshipSpec) error {
	var children []RoleRef
	for _, id := range ks.Children {
		children = append(children, RoleRef{ID: id})
	}
	kr := KinshipRoot{
		Parent:   RoleRef{ID: ks.Parent},
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
		return One{ID: id.New[ID]()}
	case TpRef:
		return TpRef{ID: id.New[ID](), Name: st.Name}
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
