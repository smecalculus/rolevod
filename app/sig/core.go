package sig

import (
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"

	"smecalculus/rolevod/app/role"
)

type ID = id.ADT
type FQN = sym.ADT
type Name = string

type Spec struct {
	// Fully Qualified Name
	FQN sym.ADT
	// Providable Endpoint Spec
	PE chnl.Spec
	// Consumable Endpoint Specs
	CEs []chnl.Spec
}

type Ref struct {
	ID   id.ADT
	Name string
}

// aka ExpDec or ExpDecDef without expression
type Root struct {
	ID       id.ADT
	Name     string
	PE       chnl.Spec
	CEs      []chnl.Spec
	Children []Ref
}

type API interface {
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (Root, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]Ref, error)
}

type service struct {
	sigs     Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newService(sigs Repo, kinships kinshipRepo, l *slog.Logger) *service {
	name := slog.String("name", "sigService")
	return &service{sigs, kinships, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	s.log.Debug("sig creation started", slog.Any("spec", spec))
	root := Root{
		ID:   id.New(),
		Name: spec.FQN.Name(),
		PE:   spec.PE,
		CEs:  spec.CEs,
	}
	err := s.sigs.Insert(root)
	if err != nil {
		s.log.Error("sig insertion failed",
			slog.Any("reason", err),
			slog.Any("sig", root),
		)
		return root, err
	}
	s.log.Debug("sig creation succeeded", slog.Any("root", root))
	return root, nil
}

func (s *service) Retrieve(rid ID) (Root, error) {
	root, err := s.sigs.SelectByID(rid)
	if err != nil {
		return Root{}, err
	}
	root.Children, err = s.sigs.SelectChildren(rid)
	if err != nil {
		return Root{}, err
	}
	return root, nil
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
	s.log.Debug("establishment succeeded", slog.Any("kinship", root))
	return nil
}

func (s *service) RetreiveAll() ([]Ref, error) {
	return s.sigs.SelectAll()
}

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectByIDs([]ID) ([]Root, error)
	SelectEnv([]ID) (map[ID]Root, error)
	SelectChildren(ID) ([]Ref, error)
}

func CollectEnv(sigs []Root) []role.ID {
	roleIDs := []role.ID{}
	for _, s := range sigs {
		roleIDs = append(roleIDs, s.PE.Role)
		for _, v := range s.CEs {
			roleIDs = append(roleIDs, v.Role)
		}
	}
	return roleIDs
}

// Kinship Relation
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
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)

func ErrRootMissingInEnv(rid ID) error {
	return fmt.Errorf("root missing in env: %v", rid)
}
