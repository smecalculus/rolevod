package deal

import (
	"log/slog"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/seat"

	"smecalculus/rolevod/app/bare/chnl"
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

type Sem interface {
	sem()
}

func (Process) sem() {}
func (Message) sem() {}
func (Service) sem() {}

// aka exec.Proc
type Process struct {
	ID  id.ADT[ID]
	Exp Step
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
	Make(TranSpec) error
}

type dealService struct {
	dealRepo    dealRepo
	kinshipRepo kinshipRepo
	log         *slog.Logger
}

func newDealService(dr dealRepo, kr kinshipRepo, l *slog.Logger) *dealService {
	name := slog.String("name", "dealService")
	return &dealService{dr, kr, l.With(name)}
}

func (s *dealService) Create(spec DealSpec) (DealRoot, error) {
	root := DealRoot{
		ID:   id.New[ID](),
		Name: spec.Name,
	}
	err := s.dealRepo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *dealService) Retrieve(id id.ADT[ID]) (DealRoot, error) {
	root, err := s.dealRepo.SelectById(id)
	if err != nil {
		return DealRoot{}, err
	}
	root.Children, err = s.dealRepo.SelectChildren(id)
	if err != nil {
		return DealRoot{}, err
	}
	return root, nil
}

func (s *dealService) RetreiveAll() ([]DealRef, error) {
	return s.dealRepo.SelectAll()
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
	err := s.kinshipRepo.Insert(root)
	if err != nil {
		return err
	}
	s.log.Debug("establishment succeed", slog.Any("kinship", root))
	return nil
}

func (c *dealService) Involve(spec PartSpec) error {
	return nil
}

func (c *dealService) Make(spec TranSpec) error {
	return nil
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
// знание о текущем шаге снаружи
type TranSpec struct {
	From Step
	To   Step
}

// знание о текущем шаге внутри
type TranSpec2 struct {
	Seat seat.SeatRef
	Cont Step
}

type TranRoot struct {
	From Step
	To   Step
}

type Label string

// aka Expression
type Step interface {
	step()
}

func (Fwd) step()    {}
func (Spawn) step()  {}
func (ExpRef) step() {}
func (Lab) step()    {}
func (Case) step()   {}
func (Send) step()   {}
func (Recv) step()   {}
func (Close) step()  {}
func (Wait) step()   {}

type Fwd struct {
	ID   id.ADT[ID]
	From chnl.Ref
	To   chnl.Ref
}

type Spawn struct {
	ID     id.ADT[ID]
	Name   string
	Comm   chnl.Ref
	Values []chnl.Ref
	Cont   Step
}

// aka ExpName
type ExpRef struct {
	ID     id.ADT[ID]
	Name   string
	Comm   chnl.Ref
	Values []chnl.Ref
}

type Lab struct {
	ID   id.ADT[ID]
	Comm chnl.Ref
	Data Label
	// Cont Step
}

type Case struct {
	ID    id.ADT[ID]
	Comm  chnl.Ref
	Conts map[Label]Step
}

type Send struct {
	ID    id.ADT[ID]
	Comm  chnl.Ref
	Value chnl.Ref
	// Cont Step
}

type Recv struct {
	ID    id.ADT[ID]
	Comm  chnl.Ref
	Value chnl.Ref
	Cont  Step
}

type Close struct {
	ID   id.ADT[ID]
	Comm chnl.Ref
}

type Wait struct {
	ID   id.ADT[ID]
	Comm chnl.Ref
	Cont Step
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
