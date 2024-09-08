package seat

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
)

type ID interface{}

type SeatSpec struct {
	Name string
}

type SeatRef struct {
	ID   id.ADT[ID]
	Name string
}

// Relation
type ChanTp struct {
	Comm  chnl.Ref
	State state.Ref
}

// Aggregate Root
// aka ExpDec or ExpDecDef without expression
type SeatRoot struct {
	ID       id.ADT[ID]
	Name     string
	Children []SeatRef
	Ctx      []ChanTp
	Zc       ChanTp
}

// Port
type SeatApi interface {
	Create(SeatSpec) (SeatRoot, error)
	Retrieve(id.ADT[ID]) (SeatRoot, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]SeatRef, error)
}

type seatService struct {
	seatRepo    seatRepo
	kinshipRepo kinshipRepo
	log         *slog.Logger
}

func newSeatService(sr seatRepo, kr kinshipRepo, l *slog.Logger) *seatService {
	name := slog.String("name", "seatService")
	return &seatService{sr, kr, l.With(name)}
}

func (s *seatService) Create(spec SeatSpec) (SeatRoot, error) {
	root := SeatRoot{
		ID:   id.New[ID](),
		Name: spec.Name,
	}
	err := s.seatRepo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *seatService) Retrieve(id id.ADT[ID]) (SeatRoot, error) {
	root, err := s.seatRepo.SelectById(id)
	if err != nil {
		return SeatRoot{}, err
	}
	root.Children, err = s.seatRepo.SelectChildren(id)
	if err != nil {
		return SeatRoot{}, err
	}
	return root, nil
}

func (s *seatService) Establish(spec KinshipSpec) error {
	var children []SeatRef
	for _, id := range spec.ChildrenIDs {
		children = append(children, SeatRef{ID: id})
	}
	root := KinshipRoot{
		Parent:   SeatRef{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinshipRepo.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", root))
	return nil
}

func (s *seatService) RetreiveAll() ([]SeatRef, error) {
	return s.seatRepo.SelectAll()
}

// Port
type seatRepo interface {
	Insert(SeatRoot) error
	SelectAll() ([]SeatRef, error)
	SelectById(id.ADT[ID]) (SeatRoot, error)
	SelectChildren(id.ADT[ID]) ([]SeatRef, error)
}

type KinshipSpec struct {
	ParentID    id.ADT[ID]
	ChildrenIDs []id.ADT[ID]
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
// goverter:extend to.*
var (
	ToSeatRef func(SeatRoot) SeatRef
	ToCore    func([]string) ([]id.ADT[ID], error)
	ToEdge    func([]id.ADT[ID]) []string
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
