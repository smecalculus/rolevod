package sig

import (
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
)

type ID = id.ADT
type FQN = sym.ADT
type Name = string

type Spec struct {
	// Fully Qualified Name
	FQN FQN
	// Providable Endpoint Spec
	PE chnl.Spec
	// Consumable Endpoint Specs
	CEs []chnl.Spec
}

type Ref struct {
	ID ID
	// Short Name
	Name Name
}

// aka ExpDec or ExpDecDef without expression
type Root struct {
	ID       ID
	Name     Name
	PE       chnl.Spec
	CEs      []chnl.Spec
	Children []Ref
}

type Api interface {
	Create(Spec) (Root, error)
	Retrieve(ID) (Root, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]Ref, error)
}

type sigService struct {
	sigs     Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newSigService(sigs Repo, kinships kinshipRepo, l *slog.Logger) *sigService {
	name := slog.String("name", "sigService")
	return &sigService{sigs, kinships, l.With(name)}
}

func (s *sigService) Create(spec Spec) (Root, error) {
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

func (s *sigService) Retrieve(rid ID) (Root, error) {
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

func (s *sigService) Establish(spec KinshipSpec) error {
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

func (s *sigService) RetreiveAll() ([]Ref, error) {
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

func CollectStIDs(sigs []Root) []state.ID {
	stIDs := []state.ID{}
	for _, s := range sigs {
		stIDs = append(stIDs, s.PE.StID)
		for _, v := range s.CEs {
			stIDs = append(stIDs, v.StID)
		}
	}
	return stIDs
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
// goverter:extend smecalculus/rolevod/lib/id:Ident
var (
	ConvertRootToRef func(Root) Ref
)

func ErrRootMissingInEnv(rid ID) error {
	return fmt.Errorf("root missing in env: %v", rid)
}
