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
	procs    step.Repo[step.ProcRoot]
	msgs     step.Repo[step.MsgRoot]
	srvs     step.Repo[step.SrvRoot]
	states   state.Repo
	kinships kinshipRepo
	log      *slog.Logger
}

func newDealService(
	deals dealRepo,
	seats seat.SeatApi,
	chnls chnl.Repo,
	procs step.Repo[step.ProcRoot],
	msgs step.Repo[step.MsgRoot],
	srvs step.Repo[step.SrvRoot],
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
			slog.Any("root", root),
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
		ID:   id.New[chnl.ID](),
		Name: seat.Via.Name,
		St:   seat.Via.St,
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
	if spec.Term == nil {
		panic(step.ErrUnexpectedTerm(spec.Term))
	}
	s.log.Debug("transition taking started", slog.Any("spec", spec))
	switch term := spec.Term.(type) {
	case step.CloseSpec:
		curChnl, err := s.chnls.SelectByID(term.A.ID)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("channel", term.A),
			)
			return err
		}
		if curChnl.St == nil {
			return fmt.Errorf("channel already finalized %+v", curChnl)
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			return err
		}
		err = s.checkTerm(term, curSt)
		if err != nil {
			return err
		}
		// TODO выборка с проверкой потребления
		srv, err := s.srvs.SelectByCh(curChnl.ID)
		if err != nil {
			s.log.Error("service selection failed",
				slog.Any("reason", err),
				slog.Any("channel", curChnl),
			)
			return err
		}
		if srv == nil {
			newMsg := step.MsgRoot{
				ID:    id.New[step.ID](),
				ViaID: term.A.ID,
				Val:   term,
			}
			err = s.msgs.Insert(newMsg)
			if err != nil {
				s.log.Error("message insertion failed",
					slog.Any("reason", err),
					slog.Any("message", newMsg),
				)
				return err
			}
			s.log.Debug("transition taking succeeded", slog.Any("message", newMsg))
			return nil
		}
		wait, ok := srv.Cont.(step.WaitSpec)
		if !ok {
			return fmt.Errorf("unexpected cont type; want %T; got %#v", step.WaitSpec{}, wait)
		}
		// consume channel
		finChnl := chnl.Root{
			ID:    id.New[chnl.ID](),
			PreID: curChnl.ID,
			Name:  curChnl.Name,
			St:    nil,
		}
		err = s.chnls.Insert(finChnl)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("channel", finChnl),
			)
			return err
		}
		// чтобы нельзя было воспользоваться потребленным каналом
		step.Subst(wait.Cont, wait.X, chnl.ConvertRootToRef(finChnl))
		return s.Take(TranSpec{DealID: spec.DealID, Term: wait.Cont})
	case step.WaitSpec:
		curChnl, err := s.chnls.SelectByID(term.X.ID)
		if err != nil {
			s.log.Error("channel selection failed",
				slog.Any("reason", err),
				slog.Any("channel", term.X),
			)
			return err
		}
		if curChnl.St == nil {
			return fmt.Errorf("channel already finalized %+v", curChnl)
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			return err
		}
		err = s.checkTerm(term, curSt)
		if err != nil {
			return err
		}
		msg, err := s.msgs.SelectByCh(curChnl.ID)
		if err != nil {
			s.log.Error("message selection failed",
				slog.Any("reason", err),
				slog.Any("channel", curChnl),
			)
			return err
		}
		if msg == nil {
			newSrv := step.SrvRoot{
				ID:    id.New[step.ID](),
				ViaID: term.X.ID,
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
		close, ok := msg.Val.(step.CloseSpec)
		if !ok {
			return fmt.Errorf("unexpected val type; want %T; got %#v", step.CloseSpec{}, close)
		}
		// consume channel
		finChnl := chnl.Root{
			ID:    id.New[chnl.ID](),
			PreID: curChnl.ID,
			Name:  curChnl.Name,
			St:    nil,
		}
		err = s.chnls.Insert(finChnl)
		if err != nil {
			s.log.Error("channel insertion failed",
				slog.Any("reason", err),
				slog.Any("channel", finChnl),
			)
			return err
		}
		// чтобы нельзя было воспользоваться потребленным каналом
		step.Subst(term.Cont, term.X, chnl.ConvertRootToRef(finChnl))
		return s.Take(TranSpec{DealID: spec.DealID, Term: term.Cont})
	case step.SendSpec:
		curChnl, err := s.chnls.SelectByID(term.A.ID)
		if err != nil {
			return err
		}
		if curChnl.St == nil {
			return fmt.Errorf("channel already finalized %+v", curChnl)
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			return err
		}
		err = s.checkTerm(term, curSt)
		if err != nil {
			return err
		}
		srv, err := s.srvs.SelectByCh(curChnl.ID)
		if err != nil {
			return err
		}
		if srv == nil {
			newMsg := step.MsgRoot{
				ID:    id.New[step.ID](),
				ViaID: term.A.ID,
				Val:   term,
			}
			return s.msgs.Insert(newMsg)
		}
		recv, ok := srv.Cont.(step.RecvSpec)
		if !ok {
			return fmt.Errorf("unexpected cont type; want %T; got %#v", step.RecvSpec{}, recv)
		}
		// TODO смена состояния канала
		step.Subst(recv.Cont, recv.X, term.A)
		step.Subst(recv.Cont, recv.Y, term.B)
		newProc := step.ProcRoot{
			ID:    id.New[step.ID](),
			PreID: srv.ID,
			Term:  recv.Cont,
		}
		// TODO рекурсивный вызов
		return s.procs.Insert(newProc)
	case step.RecvSpec:
		curChnl, err := s.chnls.SelectByID(term.X.ID)
		if err != nil {
			return err
		}
		if curChnl.St == nil {
			return fmt.Errorf("channel already finalized %+v", curChnl)
		}
		curSt, err := s.states.SelectByID(curChnl.St.RID())
		if err != nil {
			return err
		}
		err = s.checkTerm(term, curSt)
		if err != nil {
			return err
		}
		msg, err := s.msgs.SelectByCh(curChnl.ID)
		if err != nil {
			return err
		}
		if msg == nil {
			newSrv := step.SrvRoot{
				ID:    id.New[step.ID](),
				ViaID: term.X.ID,
				Cont:  term,
			}
			return s.srvs.Insert(newSrv)
		}
		send, ok := msg.Val.(step.SendSpec)
		if !ok {
			return fmt.Errorf("unexpected val type; want %T; got %#v", step.SendSpec{}, send)
		}
		step.Subst(term.Cont, term.X, send.A)
		step.Subst(term.Cont, term.Y, send.B)
		newProc := step.ProcRoot{
			ID:    id.New[step.ID](),
			PreID: msg.ID,
			Term:  term.Cont,
		}
		return s.procs.Insert(newProc)
	default:
		panic(step.ErrUnexpectedTerm(spec.Term))
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

// aka checkExp
func (s *dealService) checkTerm(t step.Term, c state.Spec) error {
	switch term := t.(type) {
	case step.CloseSpec:
		return checkState(c, state.OneSpec{})
	case step.WaitSpec:
		return checkState(c, state.OneSpec{})
	case step.SendSpec:
		// check value
		st, ok := c.(state.TensorSpec)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", state.TensorSpec{}, st)
		}
		val, err := s.chnls.SelectByID(term.B.ID)
		if err != nil {
			return err
		}
		valSt, err := s.states.SelectByID(val.St.RID())
		if err != nil {
			return err
		}
		err = checkState(valSt, st.A)
		if err != nil {
			return err
		}
		// no cont to check
		return nil
	case step.RecvSpec:
		// check value
		st, ok := c.(state.LolliSpec)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", state.LolliSpec{}, st)
		}
		val, err := s.chnls.SelectByID(term.Y.ID)
		if err != nil {
			return err
		}
		valSt, err := s.states.SelectByID(val.St.RID())
		if err != nil {
			return err
		}
		err = checkState(valSt, st.X)
		if err != nil {
			return err
		}
		// check cont
		return s.checkTerm(term.Cont, st.Z)
	default:
		panic(step.ErrUnexpectedTerm(t))
	}
}

// aka eqtp
func checkState(g state.Spec, w state.Spec) error {
	switch want := w.(type) {
	case state.OneSpec:
		got, ok := g.(state.OneSpec)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		return nil
	case state.TensorSpec:
		got, ok := g.(state.TensorSpec)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		err := checkState(got.A, want.A)
		if err != nil {
			return err
		}
		return checkState(got.C, want.C)
	case state.LolliSpec:
		got, ok := g.(state.LolliSpec)
		if !ok {
			return fmt.Errorf("state mismatch: want %T, got %#v", want, got)
		}
		err := checkState(got.X, want.X)
		if err != nil {
			return err
		}
		return checkState(got.Z, want.Z)
	default:
		panic(state.ErrUnexpectedSpec(g))
	}
}
