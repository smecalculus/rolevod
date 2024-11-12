package role

import (
	"fmt"
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

type Snap struct {
	ID    id.ADT
	Rev   rev.ADT
	State state.Spec
	// Parts   []Ref
}

// aka TpDef
type Root struct {
	ID   id.ADT
	Rev  rev.ADT
	Name string
	// specification relation
	StateID state.ID
	// composition relation
	WholeID id.ADT
}

type API interface {
	Incept(sym.ADT) (Ref, error)
	Create(Spec) (Snap, error)
	Modify(Snap) (Snap, error)
	Retrieve(id.ADT) (Snap, error)
	RetrieveRoot(id.ADT) (Root, error)
	RetrieveSnap(Root) (Snap, error)
	RetreiveRefs() ([]Ref, error)
}

type service struct {
	roles   Repo
	states  state.Repo
	aliases alias.Repo
	log     *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func newService(
	roles Repo,
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
			slog.Any("root", newAlias),
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

func (s *service) Create(spec Spec) (Snap, error) {
	s.log.Debug("role creation started", slog.Any("spec", spec))
	newAlias := alias.Root{Sym: spec.FQN, ID: id.New(), Rev: rev.New()}
	err := s.aliases.Insert(newAlias)
	if err != nil {
		s.log.Error("alias insertion failed",
			slog.Any("reason", err),
			slog.Any("root", newAlias),
		)
		return Snap{}, err
	}
	newState := state.ConvertSpecToRoot(spec.State)
	err = s.states.Insert(newState)
	if err != nil {
		s.log.Error("state insertion failed",
			slog.Any("reason", err),
			slog.Any("root", newState),
		)
		return Snap{}, err
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
		return Snap{}, err
	}
	s.log.Debug("role creation succeeded", slog.Any("root", newRoot))
	return Snap{
		ID:    newRoot.ID,
		Rev:   newRoot.Rev,
		State: state.ConvertRootToSpec(newState),
		// State: newState,
	}, nil
}

func (s *service) Modify(newSnap Snap) (Snap, error) {
	s.log.Debug("role modification started", slog.Any("snap", newSnap))
	root, err := s.roles.SelectByID(newSnap.ID)
	if err != nil {
		s.log.Error("root selection failed",
			slog.Any("reason", err),
			slog.Any("id", newSnap.ID),
		)
		return Snap{}, err
	}
	if newSnap.Rev != root.Rev {
		err := errConcurrentModification(newSnap.Rev, root.Rev)
		s.log.Error("role modification failed",
			slog.Any("reason", err),
			slog.Any("id", root.ID),
		)
		return Snap{}, err
	} else {
		newSnap.Rev = rev.New()
	}
	curSnap, err := s.RetrieveSnap(root)
	if err != nil {
		s.log.Error("snapshot retrieval failed",
			slog.Any("reason", err),
			slog.Any("root", root),
		)
		return Snap{}, err
	}
	diff := state.CheckSpec(newSnap.State, curSnap.State)
	if diff != nil {
		newState := state.ConvertSpecToRoot(newSnap.State)
		err := s.states.Insert(newState)
		if err != nil {
			s.log.Error("state insertion failed",
				slog.Any("reason", err),
				slog.Any("root", newState),
			)
			return Snap{}, err
		}
		root.StateID = newState.Ident()
		root.Rev = newSnap.Rev
	}
	if root.Rev == newSnap.Rev {
		err := s.roles.Insert(root)
		if err != nil {
			s.log.Error("role insertion failed",
				slog.Any("reason", err),
				slog.Any("root", root),
			)
			return Snap{}, err
		}
	}
	s.log.Debug("role modification succeeded", slog.Any("root", root))
	return newSnap, nil
}

func (s *service) Retrieve(eid ID) (Snap, error) {
	return Snap{}, nil
}

func (s *service) RetrieveRoot(eid ID) (Root, error) {
	root, err := s.roles.SelectByID(eid)
	if err != nil {
		s.log.Error("root selection failed")
		return Root{}, err
	}
	return root, nil
}

func (s *service) RetrieveSnap(root Root) (Snap, error) {
	snapState, err := s.states.SelectByID(root.StateID)
	if err != nil {
		s.log.Error("state selection failed")
		return Snap{}, err
	}
	return Snap{
		ID:    root.ID,
		Rev:   root.Rev,
		State: state.ConvertRootToSpec(snapState),
	}, nil
}

func (s *service) RetreiveRefs() ([]Ref, error) {
	return s.roles.SelectAll()
}

func CollectEnv(roles []Root) []id.ADT {
	roleIDs := []id.ADT{}
	for _, r := range roles {
		roleIDs = append(roleIDs, r.StateID)
	}
	return roleIDs
}

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(id.ADT) (Root, error)
	SelectParts(id.ADT) ([]Ref, error)
	SelectByIDs([]ID) ([]Root, error)
	SelectEnv([]ID) (map[ID]Root, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Convert.*
var (
	ConvertRootToRef func(Root) Ref
	// goverter:ignore Name
	ConvertSnapToRef func(Snap) Ref
)

func errConcurrentModification(got rev.ADT, want rev.ADT) error {
	return fmt.Errorf("entity concurrent modification: want revision %v, got revision %v", want, got)
}
