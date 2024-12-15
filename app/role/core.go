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
type Title = string

type Spec struct {
	FQN   sym.ADT
	State state.Spec
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
	FQN   sym.ADT
	State state.Spec
	// Parts   []Ref
}

// aka TpDef
type Root struct {
	ID    id.ADT
	Rev   rev.ADT
	Title string
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
	err = s.roles.Insert(newRoot)
	if err != nil {
		s.log.Error("role insertion failed",
			slog.Any("reason", err),
			slog.Any("root", newRoot),
		)
		return Ref{}, err
	}
	s.log.Debug("role inception succeeded", slog.Any("id", newRoot.ID))
	return ConvertRootToRef(newRoot), nil
}

func (s *service) Create(spec Spec) (Snap, error) {
	s.log.Debug("role creation started", slog.Any("spec", spec))
	newAlias := alias.Root{Sym: spec.FQN, ID: id.New(), Rev: rev.Initial()}
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
		Title:   newAlias.Sym.Name(),
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
	s.log.Debug("role creation succeeded", slog.Any("id", newRoot.ID))
	return Snap{
		ID:    newRoot.ID,
		Rev:   newRoot.Rev,
		Title: newRoot.Title,
		FQN:   newAlias.Sym,
		State: state.ConvertRootToSpec(newState),
	}, nil
}

func (s *service) Modify(newSnap Snap) (Snap, error) {
	s.log.Debug("role modification started", slog.Any("snap", newSnap))
	curRoot, err := s.roles.SelectByID(newSnap.ID)
	if err != nil {
		s.log.Error("root selection failed",
			slog.Any("reason", err),
			slog.Any("id", newSnap.ID),
		)
		return Snap{}, err
	}
	if newSnap.Rev != curRoot.Rev {
		err := errConcurrentModification(newSnap.Rev, curRoot.Rev)
		s.log.Error("role modification failed",
			slog.Any("reason", err),
			slog.Any("id", curRoot.ID),
		)
		return Snap{}, err
	} else {
		newSnap.Rev = rev.Next(newSnap.Rev)
	}
	curSnap, err := s.RetrieveSnap(curRoot)
	if err != nil {
		s.log.Error("snapshot retrieval failed",
			slog.Any("reason", err),
			slog.Any("root", curRoot),
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
		curRoot.StateID = newState.Ident()
		curRoot.Rev = newSnap.Rev
	}
	if curRoot.Rev == newSnap.Rev {
		err := s.roles.Update(curRoot)
		if err != nil {
			s.log.Error("root update failed",
				slog.Any("reason", err),
				slog.Any("root", curRoot),
			)
			return Snap{}, err
		}
	}
	s.log.Debug("role modification succeeded", slog.Any("id", curRoot.ID))
	return newSnap, nil
}

func (s *service) Retrieve(rid ID) (Snap, error) {
	root, err := s.roles.SelectByID(rid)
	if err != nil {
		s.log.Error("root selection failed", slog.Any("reason", err))
		return Snap{}, err
	}
	return s.RetrieveSnap(root)
}

func (s *service) RetrieveRoot(rid ID) (Root, error) {
	root, err := s.roles.SelectByID(rid)
	if err != nil {
		s.log.Error("root selection failed", slog.Any("reason", err))
		return Root{}, err
	}
	return root, nil
}

func (s *service) RetrieveSnap(root Root) (Snap, error) {
	curState, err := s.states.SelectByID(root.StateID)
	if err != nil {
		s.log.Error("state selection failed", slog.Any("reason", err))
		return Snap{}, err
	}
	return Snap{
		ID:    root.ID,
		Rev:   root.Rev,
		Title: root.Title,
		State: state.ConvertRootToSpec(curState),
	}, nil
}

func (s *service) RetreiveRefs() ([]Ref, error) {
	return s.roles.SelectRefs()
}

func CollectEnv(roles []Root) []id.ADT {
	stateIDs := []id.ADT{}
	for _, r := range roles {
		stateIDs = append(stateIDs, r.StateID)
	}
	return stateIDs
}

type Repo interface {
	Insert(Root) error
	Update(Root) error
	SelectRefs() ([]Ref, error)
	SelectByID(id.ADT) (Root, error)
	SelectByIDs([]id.ADT) ([]Root, error)
	SelectByFQN(sym.ADT) (Root, error)
	SelectByFQNs([]sym.ADT) ([]Root, error)
	// SelectByRef(Ref) (Snap, error)
	SelectParts(id.ADT) ([]Ref, error)
	SelectEnv([]sym.ADT) (map[sym.ADT]Root, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Convert.*
var (
	ConvertRootToRef func(Root) Ref
	ConvertSnapToRef func(Snap) Ref
)

func errConcurrentModification(got rev.ADT, want rev.ADT) error {
	return fmt.Errorf("entity concurrent modification: want revision %v, got revision %v", want, got)
}

func errOptimisticUpdate(got rev.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
