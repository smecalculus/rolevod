package seat

import (
	"errors"
	"log/slog"

	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/lib/core"
)

var (
	ErrUnexpectedExp = errors.New("unexpected expression type")
)

type Expname = string

// Aggregate Root
type Seat interface {
	seat()
}

type SeatSpec struct {
	Name Expname
}

type SeatTeaser struct {
	ID   core.ID[Seat]
	Name Expname
}

type Label string
type Chan struct {
	V string
}
type ChanTp struct {
	X  Chan
	Tp role.Stype
}

type Context []ChanTp
type Branches map[Label]Expression

// aka ExpDec or ExpDecDef without expression
type SeatRoot struct {
	ID       core.ID[Seat]
	Name     Expname
	Ctx      Context
	Zc       ChanTp
	Children []SeatTeaser
}

func (SeatRoot) seat() {}

type Expression interface {
	exp()
}

func (Fwd) exp()    {}
func (Spawn) exp()  {}
func (ExpRef) exp() {}
func (Lab) exp()    {}
func (Case) exp()   {}
func (Send) exp()   {}
func (Recv) exp()   {}
func (Close) exp()  {}
func (Wait) exp()   {}

type Fwd struct {
	ID core.ID[Seat]
	X  Chan
	Y  Chan
}

type Spawn struct {
	ID   core.ID[Seat]
	Name Expname
	Xs   []Chan
	X    Chan
	Q    Expression
}

// aka ExpName
type ExpRef struct {
	ID   core.ID[Seat]
	Name Expname
	Xs   []Chan
	X    Chan
}

type Lab struct {
	ID  core.ID[Seat]
	Ch  Chan
	L   Label
	Exp Expression
}

type Case struct {
	ID  core.ID[Seat]
	Ch  Chan
	Brs Branches
}

type Send struct {
	ID  core.ID[Seat]
	Ch1 Chan
	Ch2 Chan
	Exp Expression
}

type Recv struct {
	ID  core.ID[Seat]
	Ch1 Chan
	Ch2 Chan
	Exp Expression
}

type Close struct {
	ID core.ID[Seat]
	X  Chan
}

type Wait struct {
	ID core.ID[Seat]
	X  Chan
	P  Expression
}

// Port
type SeatApi interface {
	Create(SeatSpec) (SeatRoot, error)
	Retrieve(core.ID[Seat]) (SeatRoot, error)
	Establish(KinshipSpec) error
	RetreiveAll() ([]SeatTeaser, error)
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
		ID:   core.New[Seat](),
		Name: spec.Name,
	}
	err := s.seatRepo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *seatService) Retrieve(id core.ID[Seat]) (SeatRoot, error) {
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
	var children []SeatTeaser
	for _, id := range spec.ChildrenIDs {
		children = append(children, SeatTeaser{ID: id})
	}
	root := KinshipRoot{
		Parent:   SeatTeaser{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinshipRepo.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", root))
	return nil
}

func (s *seatService) RetreiveAll() ([]SeatTeaser, error) {
	return s.seatRepo.SelectAll()
}

// Port
type seatRepo interface {
	Insert(SeatRoot) error
	SelectById(core.ID[Seat]) (SeatRoot, error)
	SelectChildren(core.ID[Seat]) ([]SeatTeaser, error)
	SelectAll() ([]SeatTeaser, error)
}

type KinshipSpec struct {
	ParentID    core.ID[Seat]
	ChildrenIDs []core.ID[Seat]
}

type KinshipRoot struct {
	Parent   SeatTeaser
	Children []SeatTeaser
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToSeatTeaser func(SeatRoot) SeatTeaser
)

func toSame(id core.ID[Seat]) core.ID[Seat] {
	return id
}

func toCore(id string) (core.ID[Seat], error) {
	return core.FromString[Seat](id)
}

func toEdge(id core.ID[Seat]) string {
	return id.String()
}
