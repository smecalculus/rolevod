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
	Name string
}

type DealRef struct {
	ID   id.ADT[ID]
	Name string
}

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
	ConvertRootToRef func(DealRoot) DealRef
)

type DealApi interface {
	Create(DealSpec) (DealRoot, error)
	Retrieve(id.ADT[ID]) (DealRoot, error)
	RetreiveAll() ([]DealRef, error)
	Establish(KinshipSpec) error
	Involve(PartSpec) (chnl.Ref, error)
	Take(TranSpec) error
}

type dealService struct {
	deals    dealRepo
	seats    seat.SeatApi
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
	seats seat.SeatApi,
	chnls chnl.Repo,
	procs step.Repo[step.Process],
	msgs step.Repo[step.Message],
	srvs step.Repo[step.Service],
	states state.Repo,
	kinships kinshipRepo,
	l *slog.Logger,
) *dealService {
	name := slog.String("name", "dealService")
	return &dealService{deals, seats, chnls, procs, msgs, srvs, states, kinships, l.With(name)}
}

func (s *dealService) Create(spec DealSpec) (DealRoot, error) {
	s.log.Debug("deal creation started", slog.Any("spec", spec))
	root := DealRoot{
		ID:   id.New[ID](),
		Name: spec.Name,
	}
	err := s.deals.Insert(root)
	if err != nil {
		s.log.Error("deal insertion failed",
			slog.Any("reason", err),
			slog.Any("deal", root),
		)
		return root, err
	}
	s.log.Debug("deal creation succeeded", slog.Any("root", root))
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

func (s *dealService) Establish(spec KinshipSpec) error {
	s.log.Debug("kinship establishment started", slog.Any("spec", spec))
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
	s.log.Debug("kinship establishment succeeded", slog.Any("root", root))
	return nil
}

func (s *dealService) Involve(spec PartSpec) (chnl.Ref, error) {
	s.log.Debug("seat involvement started", slog.Any("spec", spec))
	seat, err := s.seats.Retrieve(spec.SeatID)
	if err != nil {
		s.log.Error("seat selection failed",
			slog.Any("reason", err),
			slog.Any("spec", spec),
		)
		return chnl.Ref{}, err
	}
	// производим внешний spawn
	// TODO переоформление контекста
	newChnl := chnl.Root{
		ID:    id.New[chnl.ID](),
		Name:  string(seat.Via.Z),
		State: seat.Via.State,
	}
	err = s.chnls.Insert(newChnl)
	if err != nil {
		s.log.Error("channel insertion failed",
			slog.Any("reason", err),
			slog.Any("channel", newChnl),
		)
		return chnl.Ref{}, err
	}
	s.log.Debug("seat involvement succeeded", slog.Any("channel", newChnl))
	return chnl.ConvertRootToRef(newChnl), nil
}

func (s *dealService) Take(spec TranSpec) error {
	s.log.Debug("transition taking started", slog.Any("spec", spec))
	switch term := spec.Term.(type) {
	case step.Close:
		curChnl, err := s.chnls.SelectByID(term.A.ID)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("channel", term.A),
			)
			return err
		}
		// TODO typecheck
		switch curChnl.State.(type) {
		case nil:
			return fmt.Errorf("channel already finalized %+v", curChnl)
		case state.OneRef:
			// TODO выборка с проверкой потребления
			srv, err := s.srvs.SelectByChID(curChnl.ID)
			if err != nil {
				s.log.Error("service selection failed",
					slog.Any("reason", err),
					slog.Any("channel", curChnl),
				)
				return err
			}
			if srv == nil {
				newMsg := step.Message{
					ID:    id.New[step.ID](),
					ViaID: term.A.ID,
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
				PreID: curChnl.ID,
				State: nil,
			}
			err = s.chnls.Insert(finChnl)
			if err != nil {
				s.log.Error("channel insertion failed",
					slog.Any("reason", err),
					slog.Any("channel", finChnl),
				)
				return err
			}
			// consume service
			finSrv := step.Service{
				PreID: srv.ID,
				Cont:  nil,
			}
			return s.srvs.Insert(finSrv)
		default:
			return fmt.Errorf("unexpected channel state; want %T; got %#v", state.OneRoot{}, curChnl)
		}
	case step.Wait:
		curChnl, err := s.chnls.SelectByID(term.X.ID)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("channel", term.X),
			)
			return err
		}
		switch curState := curChnl.State.(type) {
		case nil:
			return fmt.Errorf("channel already finalized %+v", curChnl)
		case state.OneRef:
			msg, err := s.msgs.SelectByChID(curChnl.ID)
			if err != nil {
				s.log.Error("message selection failed",
					slog.Any("reason", err),
					slog.Any("channel", curChnl),
				)
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
					s.log.Error("channel insertion failed",
						slog.Any("reason", err),
						slog.Any("channel", newChnl),
					)
					return err
				}
				viaID := term.X.ID
				term.X = chnl.ConvertRootToRef(newChnl)
				newSrv := step.Service{
					ID:    id.New[step.ID](),
					ViaID: viaID,
					Cont:  term,
				}
				err = s.srvs.Insert(newSrv)
				if err != nil {
					s.log.Error("service insertion failed",
						slog.Any("reason", err),
						slog.Any("service", newSrv),
					)
					return err
				}
				s.log.Debug("transition taking succeeded", slog.Any("service", newSrv))
				return nil
			}
			close, ok := msg.Val.(step.Close)
			if !ok {
				return fmt.Errorf("unexpected val type; want %T; got %#v", step.Close{}, close)
			}
			// consume channel
			finChnl := chnl.Root{
				PreID: curChnl.ID,
				State: nil,
			}
			err = s.chnls.Insert(finChnl)
			if err != nil {
				s.log.Error("channel insertion failed",
					slog.Any("reason", err),
					slog.Any("channel", finChnl),
				)
				return err
			}
			// consume message
			finMsg := step.Message{
				PreID: msg.ID,
				Val:   nil,
			}
			err = s.msgs.Insert(finMsg)
			if err != nil {
				s.log.Error("message insertion failed",
					slog.Any("reason", err),
					slog.Any("message", finMsg),
				)
				return err
			}
			s.log.Debug("transition taking succeeded", slog.Any("message", finMsg))
			return nil
		default:
			return fmt.Errorf("unexpected channel state; want %T; got %#v", state.OneRoot{}, curChnl)
		}
	case step.Send:
		curChnl, err := s.chnls.SelectByID(term.A.ID)
		if err != nil {
			return err
		}
		switch curState := curChnl.State.(type) {
		case nil:
			return fmt.Errorf("channel already finalized %+v", curChnl)
		case state.TensorRef:
			srv, err := s.srvs.SelectByChID(curChnl.ID)
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
				viaID := term.A.ID
				term.A = chnl.ConvertRootToRef(nextChnl)
				newMsg := step.Message{
					ID:    id.New[step.ID](),
					ViaID: viaID,
					Val:   term,
				}
				return s.msgs.Insert(newMsg)
			}
			recv, ok := srv.Cont.(step.Recv)
			if !ok {
				return fmt.Errorf("unexpected cont type; want %T; got %#v", step.Recv{}, recv)
			}
			// TODO смена состояния канала
			step.Subst(recv.Cont, recv.X, term.A)
			step.Subst(recv.Cont, recv.Y, term.B)
			newProc := step.Process{
				ID:    id.New[step.ID](),
				PreID: srv.ID,
				Term:  recv.Cont,
			}
			// TODO рекурсивный вызов
			return s.procs.Insert(newProc)
		default:
			return fmt.Errorf("unexpected channel state; want %T; got %#v", state.TensorRoot{}, curChnl)
		}
	case step.Recv:
		curChnl, err := s.chnls.SelectByID(term.X.ID)
		if err != nil {
			return err
		}
		switch curState := curChnl.State.(type) {
		case nil:
			return fmt.Errorf("channel already finalized %+v", curChnl)
		case state.LolliRef:
			msg, err := s.msgs.SelectByChID(curChnl.ID)
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
				viaID := term.X.ID
				term.X = chnl.ConvertRootToRef(newChnl)
				newSrv := step.Service{
					ID:    id.New[step.ID](),
					ViaID: viaID,
					Cont:  term,
				}
				return s.srvs.Insert(newSrv)
			}
			send, ok := msg.Val.(step.Send)
			if !ok {
				return fmt.Errorf("unexpected val type; want %T; got %#v", step.Send{}, send)
			}
			step.Subst(term.Cont, term.X, send.A)
			step.Subst(term.Cont, term.Y, send.B)
			newProc := step.Process{
				ID:    id.New[step.ID](),
				PreID: msg.ID,
				Term:  term.Cont,
			}
			return s.procs.Insert(newProc)
		default:
			return fmt.Errorf("unexpected channel state; want %T; got %#v", state.TensorRoot{}, curChnl)
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
	DealID id.ADT[ID]
	SeatID id.ADT[seat.ID]
}

type PartRoot struct {
	Deal DealRef
	Seat seat.SeatRef
}

type partRepo interface {
	Insert(PartRoot) error
}

// Transition
type TranSpec struct {
	DealID id.ADT[ID]
	Term   step.Term
}

type TranRoot struct {
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
