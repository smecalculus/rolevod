package deal

import (
	"errors"
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/seat"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"
)

type ID interface{}

type DealSpec struct {
	Name  string
	Seats []seat.SeatRef
}

type DealRef struct {
	ID   id.ADT[ID]
	Name string
}

// Aggregate Root
// aka Configuration or Eta
type DealRoot struct {
	ID       id.ADT[ID]
	Name     string
	Children []DealRef
	Seats    []seat.SeatRef
	Prcs     map[chnl.Ref]Process
	Msgs     map[chnl.Ref]Message
	Srvs     map[chnl.Ref]Service
}

// aka exec.Proc
type Process struct {
	ID  id.ADT[ID]
	Exp step.Root
}

// aka exec.Msg
type Message struct {
	ID      id.ADT[ID]
	Comm    chnl.Ref
	Payload Value
}

// aka ast.Msg
type Value interface {
	val()
}

type Service struct {
	ID   id.ADT[ID]
	Comm chnl.Ref
	Cont Continuation
}

type Continuation interface {
	cont()
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	ToDealRef func(DealRoot) DealRef
)

type DealApi interface {
	Create(DealSpec) (DealRoot, error)
	Retrieve(id.ADT[ID]) (DealRoot, error)
	RetreiveAll() ([]DealRef, error)
	Establish(KinshipSpec) error
	Involve(PartSpec) error
	Take(Transition) error
}

type dealService struct {
	deals    dealRepo
	chnls    chnl.Repo
	steps    step.Repo
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newDealService(
	deals dealRepo,
	chnls chnl.Repo,
	steps step.Repo,
	states state.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *dealService {
	name := slog.String("name", "dealService")
	return &dealService{deals, chnls, steps, states, kinships, l.With(name)}
}

func (s *dealService) Create(spec DealSpec) (DealRoot, error) {
	root := DealRoot{
		ID:   id.New[ID](),
		Name: spec.Name,
	}
	err := s.deals.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *dealService) Retrieve(id id.ADT[ID]) (DealRoot, error) {
	root, err := s.deals.SelectById(id)
	if err != nil {
		return DealRoot{}, err
	}
	root.Children, err = s.deals.SelectChildren(id)
	if err != nil {
		return DealRoot{}, err
	}
	return root, nil
}

func (s *dealService) RetreiveAll() ([]DealRef, error) {
	return s.deals.SelectAll()
}

func (s *dealService) Establish(spec KinshipSpec) error {
	var children []DealRef
	for _, id := range spec.ChildrenIDs {
		children = append(children, DealRef{ID: id})
	}
	root := KinshipRoot{
		Parent:   DealRef{ID: spec.ParentID},
		Children: children,
	}
	err := s.kinships.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", root))
	return nil
}

func (s *dealService) Involve(spec PartSpec) error {
	return nil
}

func (s *dealService) Take(rel Transition) error {
	switch want := rel.Step.(type) {
	case step.Close:
		curChnl, err := s.chnls.SelectById(want.Via.ID)
		if err != nil {
			return err
		}
		switch curState := curChnl.State.(type) {
		case state.One:
			got, err := s.steps.SelectByChnl(curChnl.ID)
			if err != nil {
				return err
			}
			_, ok := got.(step.Wait)
			if !ok {
				return fmt.Errorf("Unexpected continuation step %+v in channel %+v", got, curChnl)
			}
			return nil
		default:
			nextState, err := s.states.SelectNext(curState.Get())
			if err != nil {
				return err
			}
			_, ok := nextState.(state.One)
			if !ok {
				return fmt.Errorf("Unexpected next state %+v in channel %+v", nextState, curChnl)
			}
			nextChnl := chnl.Root{
				ID:    id.New[chnl.ID](),
				Name:  curChnl.Name,
				State: nextState,
			}
			want.Via = chnl.ToRef(nextChnl)
			return s.steps.Insert(want)
		}
	case step.Wait:
		curChnl, err := s.chnls.SelectById(want.Via.ID)
		if err != nil {
			return err
		}
		switch curState := curChnl.State.(type) {
		case state.One:
			got, err := s.steps.SelectByChnl(curChnl.ID)
			if err != nil {
				return err
			}
			_, ok := got.(step.Close)
			if !ok {
				return fmt.Errorf("Unexpected continuation step %+v in channel %+v", got, curChnl)
			}
			return nil
		default:
			nextState, err := s.states.SelectNext(curState.Get())
			if err != nil {
				return err
			}
			_, ok := nextState.(state.One)
			if !ok {
				return fmt.Errorf("Unexpected next state %+v in channel %+v", nextState, curChnl)
			}
			nextChnl := chnl.Root{
				ID:    id.New[chnl.ID](),
				Name:  curChnl.Name,
				State: nextState,
			}
			want.Via = chnl.ToRef(nextChnl)
			return s.steps.Insert(want)
		}
	case step.Send:
		curChnl, err := s.chnls.SelectById(want.Via.ID)
		if err != nil {
			return err
		}
		switch curState := curChnl.State.(type) {
		case state.Tensor:
			got, err := s.steps.SelectByChnl(curChnl.ID)
			if err != nil {
				return err
			}
			_, ok := got.(step.Recv)
			if !ok {
				return fmt.Errorf("Unexpected continuation step %+v in channel %+v", got, curChnl)
			}
			return nil
		default:
			nextState, err := s.states.SelectNext(curState.Get())
			if err != nil {
				return err
			}
			_, ok := nextState.(state.Tensor)
			if !ok {
				return fmt.Errorf("Unexpected next state %+v in channel %+v", nextState, curChnl)
			}
			nextChnl := chnl.Root{
				ID:    id.New[chnl.ID](),
				Name:  curChnl.Name,
				State: nextState,
			}
			want.Via = chnl.ToRef(nextChnl)
			return s.steps.Insert(want)
		}
	default:
		panic(ErrUnexpectedStep)
	}
}

type dealRepo interface {
	Insert(DealRoot) error
	SelectAll() ([]DealRef, error)
	SelectById(id.ADT[ID]) (DealRoot, error)
	SelectChildren(id.ADT[ID]) ([]DealRef, error)
	SelectSeats(id.ADT[ID]) ([]seat.SeatRef, error)
}

// Kinship Relation
type KinshipSpec struct {
	ParentID    id.ADT[ID]
	ChildrenIDs []id.ADT[ID]
}

type KinshipRoot struct {
	Parent   DealRef
	Children []DealRef
}

type kinshipRepo interface {
	Insert(KinshipRoot) error
}

// Partitcipation Relation
type PartSpec struct {
	DealID  id.ADT[ID]
	SeatIDs []id.ADT[seat.ID]
}

type PartRoot struct {
	Deal  DealRef
	Seats []seat.SeatRef
}

type partRepo interface {
	Insert(PartRoot) error
}

// Transition Relation
type Transition struct {
	Deal DealRef
	Step step.Root
}

type Label string

var (
	ErrUnexpectedStep = errors.New("unexpected step type")
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
