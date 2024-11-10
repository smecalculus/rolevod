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
	FQN   sym.ADT
	State state.Spec
}

type Ref struct {
	ID   id.ADT
	Rev  rev.ADT
	Name string
}

type Patch struct {
	ID    id.ADT
	Rev   rev.ADT
	State state.Spec
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
	Update(Patch) error
	RetrieveLatest(id.ADT) (Root, error)
	Retrieve(Ref) (Root, error)
	RetreiveAll() ([]Ref, error)
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

func (s *service) Incept(fqn sym.ADT) (Ref, error) {
	s.log.Debug("role inception started", slog.Any("fqn", fqn))
	newAlias := alias.Root{Sym: fqn, ID: id.New(), Rev: rev.New()}
	err := s.aliases.Insert(newAlias)
	if err != nil {
		s.log.Error("alias insertion failed",
			slog.Any("reason", err),
			slog.Any("alias", newAlias),
		)
		return Ref{}, err
	}
	newRoot := Root{
		ID:   newAlias.ID,
		Rev:  newAlias.Rev,
		Name: newAlias.Sym.Name(),
	}
	err = s.roles.Insert(newRoot)
	if err != nil {
		s.log.Error("role insertion failed",
			slog.Any("reason", err),
			slog.Any("root", newRoot),
		)
		return Ref{}, err
	}
	s.log.Debug("role inception succeeded", slog.Any("root", newRoot))
	return ConvertRootToRef(newRoot), nil
}

func (s *service) Create(spec Spec) (Root, error) {
	s.log.Debug("role creation started", slog.Any("spec", spec))
	newAlias := alias.Root{Sym: spec.FQN, ID: id.New(), Rev: rev.New()}
	err := s.aliases.Insert(newAlias)
	if err != nil {
		s.log.Error("alias insertion failed",
			slog.Any("reason", err),
			slog.Any("alias", newAlias),
		)
		return Root{}, err
	}
	newState := state.ConvertSpecToRoot(spec.State)
	err = s.states.Insert(newState)
	if err != nil {
		s.log.Error("state insertion failed",
			slog.Any("reason", err),
			slog.Any("state", newState),
		)
		return Root{}, err
	}
	newRoot := Root{
		ID:      newAlias.ID,
		Rev:     newAlias.Rev,
		Name:    newAlias.Sym.Name(),
		StateID: newState.Ident(),
	}
	err = s.roles.Insert(newRoot)
	if err != nil {
		s.log.Error("role insertion failed",
			slog.Any("reason", err),
			slog.Any("root", newRoot),
		)
		return Root{}, err
	}
	s.log.Debug("role creation succeeded", slog.Any("root", newRoot))
	newRoot.State = state.ConvertRootToSpec(newState)
	return newRoot, nil
}

func (s *service) Update(patch Patch) (err error) {
	s.log.Debug("role update started", slog.Any("patch", patch))
	var root Root
	if patch.Rev == 0 {
		root, err = s.RetrieveLatest(patch.ID)
		if err != nil {
			s.log.Error("root retrieval failed",
				slog.Any("reason", err),
				slog.Any("id", patch.ID),
			)
			return err
		}
	} else {
		ref := Ref{ID: patch.ID, Rev: patch.Rev}
		root, err = s.Retrieve(ref)
		s.log.Error("root retrieval failed",
			slog.Any("reason", err),
			slog.Any("ref", ref),
		)
		return err
	}
	newRev := rev.New()
	err = state.CheckSpec(patch.State, root.State)
	if err != nil {
		newState := state.ConvertSpecToRoot(patch.State)
		err = s.states.Insert(newState)
		if err != nil {
			s.log.Error("state insertion failed",
				slog.Any("reason", err),
				slog.Any("state", newState),
			)
			return err
		}
		root.StateID = newState.Ident()
		root.Rev = newRev
	}
	if root.Rev == newRev {
		err := s.roles.Insert(root)
		if err != nil {
			s.log.Error("role insertion failed",
				slog.Any("reason", err),
				slog.Any("root", root),
			)
			return err
		}
	}
	s.log.Debug("role update succeeded", slog.Any("root", root))
	return nil
}

func (s *service) RetrieveLatest(rid ID) (Root, error) {
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

func (s *service) Retrieve(ref Ref) (Root, error) {
	return s.RetrieveLatest(ref.ID)
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
