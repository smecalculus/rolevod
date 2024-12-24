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
	// Providable Endpoint
	PE chnl.Spec
	// Consumable Endpoints
	CEs []chnl.Spec
}

type Ref struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
}

type Snap struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
	CEs   []chnl.Spec
	PE    chnl.Spec
}

// aka ExpDec or ExpDecDef without expression
type Root struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
	CEs   []chnl.Spec
	PE    chnl.Spec
}

type API interface {
	Incept(FQN) (Ref, error)
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (Root, error)
	RetreiveRefs() ([]Ref, error)
}

type service struct {
	sigs    Repo
	aliases alias.Repo
	log     *slog.Logger
}

func newService(sigs Repo, aliases alias.Repo, l *slog.Logger) *service {
	name := slog.String("name", "sigService")
	return &service{sigs, aliases, l.With(name)}
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
		ID:    newAlias.ID,
		Rev:   newAlias.Rev,
		Title: newAlias.Sym.Name(),
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
		ID:    id.New(),
		Rev:   rev.Initial(),
		Title: spec.FQN.Name(),
		PE:    spec.PE,
		CEs:   spec.CEs,
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

func (s *service) RetreiveRefs() ([]Ref, error) {
	return s.sigs.SelectAll()
}

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectByIDs([]ID) ([]Root, error)
	SelectEnv([]ID) (map[ID]Root, error)
}

func CollectEnv(sigs []Root) []role.FQN {
	roleFQNs := []role.FQN{}
	for _, sig := range sigs {
		roleFQNs = append(roleFQNs, sig.PE.Link)
		for _, ce := range sig.CEs {
			roleFQNs = append(roleFQNs, ce.Link)
		}
	}
	return roleFQNs
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
