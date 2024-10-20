package seat

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

type SeatSpec struct {
	// Fully Qualified Name
	FQN sym.ADT
	// Providable Endpoint Spec
	PE chnl.Spec
	// Consumable Endpoint Specs
	CEs []chnl.Spec
}

type SeatRef struct {
	ID ID
	// Short Name
	Name string
}

// aka ExpDec or ExpDecDef without expression
type SeatRoot struct {
	ID       ID
	Name     string
	PE       chnl.Spec
	CEs      []chnl.Spec
	Children []SeatRef
}

type SeatApi interface {
	Create(SeatSpec) (SeatRoot, error)
	Retrieve(ID) (SeatRoot, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]SeatRef, error)
}

type seatService struct {
	seats    SeatRepo
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newSeatService(seats SeatRepo, states state.Repo, kinships kinshipRepo, l *slog.Logger) *seatService {
	name := slog.String("name", "seatService")
	return &seatService{seats, states, kinships, l.With(name)}
}

func (s *seatService) Create(spec SeatSpec) (SeatRoot, error) {
	s.log.Debug("seat creation started", slog.Any("spec", spec))
	root := SeatRoot{
		ID:   id.New(),
		Name: spec.FQN.Name(),
		PE:   spec.PE,
		CEs:  spec.CEs,
	}
	err := s.seats.Insert(root)
	if err != nil {
		s.log.Error("seat insertion failed",
			slog.Any("reason", err),
			slog.Any("seat", root),
		)
		return root, err
	}
	s.log.Debug("seat creation succeeded", slog.Any("root", root))
	return root, nil
}

func (s *seatService) Retrieve(rid ID) (SeatRoot, error) {
	root, err := s.seats.SelectByID(rid)
	if err != nil {
		return SeatRoot{}, err
	}
	root.Children, err = s.seats.SelectChildren(rid)
	if err != nil {
		return SeatRoot{}, err
	}
	return root, nil
}

func (s *seatService) Establish(spec KinshipSpec) error {
	var children []SeatRef
	for _, id := range spec.ChildIDs {
		children = append(children, SeatRef{ID: id})
	}
	root := KinshipRoot{
		Parent:   SeatRef{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinships.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeeded", slog.Any("kinship", root))
	return nil
}

func (s *seatService) RetreiveAll() ([]SeatRef, error) {
	return s.seats.SelectAll()
}

type SeatRepo interface {
	Insert(SeatRoot) error
	SelectAll() ([]SeatRef, error)
	SelectByID(ID) (SeatRoot, error)
	SelectByIDs([]ID) ([]SeatRoot, error)
	SelectEnv([]ID) (map[ID]SeatRoot, error)
	SelectChildren(ID) ([]SeatRef, error)
}

func CollectStIDs(seats []SeatRoot) []state.ID {
	stateIDs := []state.ID{}
	for _, s := range seats {
		stateIDs = append(stateIDs, s.PE.StID)
		for _, v := range s.CEs {
			stateIDs = append(stateIDs, v.StID)
		}
	}
	return stateIDs
}

// Kinship Relation
type KinshipSpec struct {
	ParentID ID
	ChildIDs []ID
}

type KinshipRoot struct {
	Parent   SeatRef
	Children []SeatRef
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Ident
var (
	ConvertRootToRef func(SeatRoot) SeatRef
)

func ErrRootMissingInEnv(rid ID) error {
	return fmt.Errorf("root missing in env: %v", rid)
}
