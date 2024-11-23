package sig

import (
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/rev"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/alias"
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

type Snap struct {
	ID id.ADT
	// Rev  rev.ADT
	Name string
	PE   chnl.Spec
	CEs  []chnl.Spec
}

// aka ExpDec or ExpDecDef without expression
type Root struct {
	ID id.ADT
	// Rev      rev.ADT
	Name string
	PE   chnl.Spec
	CEs  []chnl.Spec
}

type API interface {
	Incept(FQN) (Ref, error)
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (Root, error)
	Establish(KinshipSpec) error
	RetreiveRefs() ([]Ref, error)
}

type service struct {
	sigs     Repo
	aliases  alias.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newService(sigs Repo, aliases alias.Repo, kinships kinshipRepo, l *slog.Logger) *service {
	name := slog.String("name", "sigService")
	return &service{sigs, aliases, kinships, l.With(name)}
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func (s *service) Incept(fqn sym.ADT) (Ref, error) {
	s.log.Debug("signature inception started", slog.Any("fqn", fqn))
	newAlias := alias.Root{Sym: fqn, ID: id.New(), Rev: rev.Initial()}
	err := s.aliases.Insert(newAlias)
	if err != nil {
		s.log.Error("alias insertion failed",
			slog.Any("reason", err),
			slog.Any("root", newAlias),
		)
		return Ref{}, err
	}
	newRoot := Root{
		ID: newAlias.ID,
		// Rev:  newAlias.Rev,
		Name: newAlias.Sym.Name(),
	}
	err = s.sigs.Insert(newRoot)
	if err != nil {
		s.log.Error("signature insertion failed",
			slog.Any("reason", err),
			slog.Any("root", newRoot),
		)
		return Ref{}, err
	}
	s.log.Debug("signature inception succeeded", slog.Any("root", newRoot))
	return ConvertRootToRef(newRoot), nil
}

func (s *service) Create(spec Spec) (Root, error) {
	s.log.Debug("signature creation started", slog.Any("spec", spec))
	root := Root{
		ID:   id.New(),
		Name: spec.FQN.Name(),
		PE:   spec.PE,
		CEs:  spec.CEs,
	}
	err := s.sigs.Insert(root)
	if err != nil {
		s.log.Error("signature insertion failed",
			slog.Any("reason", err),
			slog.Any("sig", root),
		)
		return root, err
	}
	s.log.Debug("signature creation succeeded", slog.Any("root", root))
	return root, nil
}

func (s *service) Retrieve(rid ID) (Root, error) {
	root, err := s.sigs.SelectByID(rid)
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

func (s *service) RetreiveRefs() ([]Ref, error) {
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
