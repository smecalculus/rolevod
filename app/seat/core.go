package seat

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"

	"smecalculus/rolevod/app/role"
)

type ID interface{}

type SeatSpec struct {
	Name string
	Via  ChanTp
	Ctx  []ChanTp
}

type ChanTpSpec struct {
	Z    chnl.Var
	Role role.RoleRef
}

type SeatRef struct {
	ID   id.ADT[ID]
	Name string
}

// aka ExpDec or ExpDecDef without expression
type SeatRoot struct {
	ID       id.ADT[ID]
	Name     string
	Via      ChanTp
	Ctx      []ChanTp
	Children []SeatRef
}

// TODO подобрать семантику для отношения
type ChanTp struct {
	Z    chnl.Var
	Role role.RoleRef
}

// Port
type SeatApi interface {
	Create(SeatSpec) (SeatRoot, error)
	Retrieve(id.ADT[ID]) (SeatRoot, error)
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
	root := SeatRoot{
		ID:   id.New[ID](),
		Name: spec.Name,
		Via:  spec.Via,
		Ctx:  spec.Ctx,
	}
	err := s.seats.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *seatService) Retrieve(id id.ADT[ID]) (SeatRoot, error) {
	root, err := s.seats.SelectByID(id)
	if err != nil {
		return SeatRoot{}, err
	}
	root.Children, err = s.seats.SelectChildren(id)
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
	err := s.kinships.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", root))
	return nil
}

func (s *seatService) RetreiveAll() ([]SeatRef, error) {
	return s.seats.SelectAll()
}

// Port
type SeatRepo interface {
	Insert(SeatRoot) error
	SelectAll() ([]SeatRef, error)
	SelectByID(id.ADT[ID]) (SeatRoot, error)
	SelectChildren(id.ADT[ID]) ([]SeatRef, error)
}

// Kinship Relation
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
	ToCoreIDs func([]string) ([]id.ADT[ID], error)
	ToEdgeIDs func([]id.ADT[ID]) []string
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
