package seat

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
)

type ID = id.ADT

type Sym = string

type SeatSpec struct {
	Name Sym
	Via  chnl.Spec
	Ctx  map[chnl.Name]chnl.Spec
}

type SeatRef struct {
	ID   ID
	Name string
}

// aka ExpDec or ExpDecDef without expression
type SeatRoot struct {
	ID       ID
	Name     Sym
	Via      chnl.Spec
	Ctx      map[chnl.Name]chnl.Spec
	Children []SeatRef
}

type SeatApi interface {
	Create(SeatSpec) (SeatRoot, error)
	Retrieve(id.ADT) (SeatRoot, error)
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
		Name: spec.Name,
		Via:  spec.Via,
		Ctx:  spec.Ctx,
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

func (s *seatService) Retrieve(rid id.ADT) (SeatRoot, error) {
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
	SelectByID(id.ADT) (SeatRoot, error)
	SelectChildren(id.ADT) ([]SeatRef, error)
}

// Kinship Relation
type KinshipSpec struct {
	ParentID id.ADT
	ChildIDs []id.ADT
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
