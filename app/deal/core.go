package deal

import (
	"fmt"
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/seat"
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
	procs    step.Repo[step.Process]
	msgs     step.Repo[step.Message]
	srvs     step.Repo[step.Service]
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newDealService(
	deals dealRepo,
	chnls chnl.Repo,
	procs step.Repo[step.Process],
	msgs step.Repo[step.Message],
	srvs step.Repo[step.Service],
	states state.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *dealService {
	name := slog.String("name", "dealService")
	return &dealService{deals, chnls, procs, msgs, srvs, states, kinships, l.With(name)}
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
	root, err := s.deals.SelectByID(id)
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

func (s *dealService) Establish(rel KinshipSpec) error {
	var children []DealRef
	for _, id := range rel.ChildrenIDs {
		children = append(children, DealRef{ID: id})
	}
	root := KinshipRoot{
		Parent:   DealRef{ID: rel.ParentID},
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
	switch want := rel.Term.(type) {
	case step.Close:
		curChnl, err := s.chnls.SelectByID(want.A.ID)
		if err != nil {
			return err
		}
		switch curChnl.State.(type) {
		case nil:
			return fmt.Errorf("channel already finalized %+v", curChnl)
		case state.One:
			srv, err := s.srvs.SelectByViaID(curChnl.ID)
			if err != nil {
				return err
			}
			if srv == nil {
				newMsg := step.Message{
					ID:    id.New[step.ID](),
					ViaID: want.A.ID,
					Val:   step.Unit,
				}
				return s.msgs.Insert(newMsg)
			}
			wait, ok := srv.Cont.(step.Wait)
			if !ok {
				return fmt.Errorf("unexpected cont type; want %T; got %#v", step.Wait{}, wait)
			}
			// consume channel
			finChnl := chnl.Root{
				ID:    id.New[chnl.ID](),
				PreID: curChnl.ID,
				Name:  curChnl.Name,
				State: nil,
			}
			err = s.chnls.Insert(finChnl)
			if err != nil {
				return err
			}
			// consume service
			finSrv := step.Service{
				ID:    id.New[step.ID](),
				PreID: srv.ID,
				Cont:  nil,
			}
			return s.srvs.Insert(finSrv)
		default:
			return fmt.Errorf("unexpected channel state; want %T; got %#v", state.One{}, curChnl)
		}
	case step.Wait:
		curChnl, err := s.chnls.SelectByID(want.X.ID)
		if err != nil {
			return err
		}
		switch curState := curChnl.State.(type) {
		case nil:
			return fmt.Errorf("channel already finalized %+v", curChnl)
		case state.One:
			msg, err := s.msgs.SelectByViaID(curChnl.ID)
			if err != nil {
				return err
			}
			if msg == nil {
				newChnl := chnl.Root{
					ID:    id.New[chnl.ID](),
					Name:  curChnl.Name,
					State: curState,
				}
				err = s.chnls.Insert(newChnl)
				if err != nil {
					return err
				}
				viaID := want.X.ID
				want.X = chnl.ToRef(newChnl)
				newSrv := step.Service{
					ID:    id.New[step.ID](),
					ViaID: viaID,
					Cont:  want,
				}
				return s.srvs.Insert(newSrv)
			}
			close, ok := msg.Val.(step.Close)
			if !ok {
				return fmt.Errorf("unexpected val type; want %T; got %#v", step.Close{}, close)
			}
			// consume channel
			finChnl := chnl.Root{
				ID:    id.New[chnl.ID](),
				PreID: curChnl.ID,
				Name:  curChnl.Name,
				State: nil,
			}
			err = s.chnls.Insert(finChnl)
			if err != nil {
				return err
			}
			// consume message
			finMsg := step.Message{
				ID:    id.New[step.ID](),
				PreID: msg.ID,
				Val:   nil,
			}
			return s.msgs.Insert(finMsg)
		default:
			return fmt.Errorf("unexpected channel state; want %T; got %#v", state.One{}, curChnl)
		}
	case step.Send:
		curChnl, err := s.chnls.SelectByID(want.A.ID)
		if err != nil {
			return err
		}
		switch curState := curChnl.State.(type) {
		case nil:
			return fmt.Errorf("channel already finalized %+v", curChnl)
		case state.Tensor:
			srv, err := s.srvs.SelectByViaID(curChnl.ID)
			if err != nil {
				return err
			}
			if srv == nil {
				nextChnl := chnl.Root{
					ID:    id.New[chnl.ID](),
					Name:  curChnl.Name,
					State: curState,
				}
				err = s.chnls.Insert(nextChnl)
				if err != nil {
					return err
				}
				viaID := want.A.ID
				want.A = chnl.ToRef(nextChnl)
				newMsg := step.Message{
					ID:    id.New[step.ID](),
					ViaID: viaID,
					Val:   want,
				}
				return s.msgs.Insert(newMsg)
			}
			recv, ok := srv.Cont.(step.Recv)
			if !ok {
				return fmt.Errorf("unexpected cont type; want %T; got %#v", step.Recv{}, recv)
			}
			step.Subst(recv.Cont, recv.X, want.A)
			step.Subst(recv.Cont, recv.Y, want.B)
			newProc := step.Process{
				ID:    id.New[step.ID](),
				PreID: srv.ID,
				Term:  recv.Cont,
			}
			return s.procs.Insert(newProc)
		default:
			return fmt.Errorf("unexpected channel state; want %T; got %#v", state.Tensor{}, curChnl)
		}
	case step.Recv:
		curChnl, err := s.chnls.SelectByID(want.X.ID)
		if err != nil {
			return err
		}
		switch curState := curChnl.State.(type) {
		case nil:
			return fmt.Errorf("channel already finalized %+v", curChnl)
		case state.Lolli:
			msg, err := s.msgs.SelectByViaID(curChnl.ID)
			if err != nil {
				return err
			}
			if msg == nil {
				newChnl := chnl.Root{
					ID:    id.New[chnl.ID](),
					Name:  curChnl.Name,
					State: curState,
				}
				err = s.chnls.Insert(newChnl)
				if err != nil {
					return err
				}
				viaID := want.X.ID
				want.X = chnl.ToRef(newChnl)
				newSrv := step.Service{
					ID:    id.New[step.ID](),
					ViaID: viaID,
					Cont:  want,
				}
				return s.srvs.Insert(newSrv)
			}
			send, ok := msg.Val.(step.Send)
			if !ok {
				return fmt.Errorf("unexpected val type; want %T; got %#v", step.Send{}, send)
			}
			step.Subst(want.Cont, want.X, send.A)
			step.Subst(want.Cont, want.Y, send.B)
			newProc := step.Process{
				ID:    id.New[step.ID](),
				PreID: msg.ID,
				Term:  want.Cont,
			}
			return s.procs.Insert(newProc)
		default:
			return fmt.Errorf("unexpected channel state; want %T; got %#v", state.Tensor{}, curChnl)
		}
	default:
		panic(step.ErrUnexpectedTerm)
	}
}

type dealRepo interface {
	Insert(DealRoot) error
	SelectAll() ([]DealRef, error)
	SelectByID(id.ADT[ID]) (DealRoot, error)
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
	Term step.Term
}

func toSame(id id.ADT[ID]) id.ADT[ID] {
	return id
}

func toCore(s string) (id.ADT[ID], error) {
	return id.String[ID](s)
}

func toEdge(id id.ADT[ID]) string {
	return id.String()
}
