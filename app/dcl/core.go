package dcl

import (
	"errors"
	"log/slog"

	"smecalculus/rolevod/lib/core"
)

var (
	ErrUnexpectedSt = errors.New("unexpected session type")
)

type Tpname = string
type Expname = string

type Spec interface {
	spec()
}

// Aggregate Root (aka decl)
type AR interface {
	root()
}

type TpSpec struct {
	Name Tpname
	St   Stype
}

func (TpSpec) spec() {}

type TpTeaser struct {
	ID   core.ID[AR]
	Name Tpname
}

func (TpRoot) root() {}

// aka TpDef
type TpRoot struct {
	ID   core.ID[AR]
	Name Tpname
	St   Stype
}

type Label string
type Choices map[Label]Stype

type Chan struct {
	V string
}

type Stype interface {
	stype()
}

func (One) stype()    {}
func (TpRef) stype()  {}
func (Tensor) stype() {}
func (Lolli) stype()  {}
func (With) stype()   {}
func (Plus) stype()   {}
func (Up) stype()     {}
func (Down) stype()   {}

// External Choice
type With struct {
	ID  core.ID[AR]
	Chs Choices
}

// Internal Choice
type Plus struct {
	ID  core.ID[AR]
	Chs Choices
}

type Tensor struct {
	ID core.ID[AR]
	S  Stype
	T  Stype
}

type Lolli struct {
	ID core.ID[AR]
	S  Stype
	T  Stype
}

type One struct {
	ID core.ID[AR]
}

// aka TpName
type TpRef struct {
	ID   core.ID[AR]
	Name Tpname
}

type Up struct {
	ID core.ID[AR]
	A  Stype
}

type Down struct {
	ID core.ID[AR]
	A  Stype
}

type ChanTp struct {
	X  Chan
	Tp Stype
}

// Port
type TpApi interface {
	Create(TpSpec) (TpRoot, error)
	Update(TpRoot) error
	Retrieve(core.ID[AR]) (TpRoot, error)
	RetreiveAll() ([]TpRoot, error)
}

type tpService struct {
	repo repo[TpRoot]
	log  *slog.Logger
}

func newTpService(r repo[TpRoot], l *slog.Logger) *tpService {
	name := slog.String("name", "dcl.tpService")
	return &tpService{r, l.With(name)}
}

func (s *tpService) Create(spec TpSpec) (TpRoot, error) {
	root := TpRoot{
		ID:   core.New[AR](),
		Name: spec.Name,
		St:   elab(spec.St),
	}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *tpService) Update(root TpRoot) error {
	return s.repo.Insert(root)
}

func (s *tpService) Retrieve(id core.ID[AR]) (TpRoot, error) {
	return s.repo.SelectById(id)
}

func (s *tpService) RetreiveAll() ([]TpRoot, error) {
	return s.repo.SelectAll()
}

type ExpSpec struct {
	Name Expname
}

func (ExpSpec) spec() {}

type ExpTeaser struct {
	ID   core.ID[AR]
	Name Expname
}

type Context []ChanTp
type Branches map[Label]Expression

// aka ExpDec or ExpDecDef without expression
type ExpRoot struct {
	ID   core.ID[AR]
	Name Expname
	Ctx  Context
	Zc   ChanTp
}

func (ExpRoot) root() {}

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
	ID core.ID[AR]
	X  Chan
	Y  Chan
}

type Spawn struct {
	ID   core.ID[AR]
	Name Expname
	Xs   []Chan
	X    Chan
	Q    Expression
}

// aka ExpName
type ExpRef struct {
	ID   core.ID[AR]
	Name Expname
	Xs   []Chan
	X    Chan
}

type Lab struct {
	ID  core.ID[AR]
	Ch  Chan
	L   Label
	Exp Expression
}

type Case struct {
	ID  core.ID[AR]
	Ch  Chan
	Brs Branches
}

type Send struct {
	ID  core.ID[AR]
	Ch1 Chan
	Ch2 Chan
	Exp Expression
}

type Recv struct {
	ID  core.ID[AR]
	Ch1 Chan
	Ch2 Chan
	Exp Expression
}

type Close struct {
	ID core.ID[AR]
	X  Chan
}

type Wait struct {
	ID core.ID[AR]
	X  Chan
	P  Expression
}

// Port
type ExpApi interface {
	Create(ExpSpec) (ExpRoot, error)
	Retrieve(core.ID[AR]) (ExpRoot, error)
	RetreiveAll() ([]ExpRoot, error)
}

type expService struct {
	repo repo[ExpRoot]
	log  *slog.Logger
}

func newExpService(r repo[ExpRoot], l *slog.Logger) *expService {
	name := slog.String("name", "dcl.expService")
	return &expService{r, l.With(name)}
}

func (s *expService) Create(spec ExpSpec) (ExpRoot, error) {
	root := ExpRoot{
		ID:   core.New[AR](),
		Name: spec.Name,
	}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *expService) Retrieve(id core.ID[AR]) (ExpRoot, error) {
	return s.repo.SelectById(id)
}

func (s *expService) RetreiveAll() ([]ExpRoot, error) {
	return s.repo.SelectAll()
}

func elab(stype Stype) Stype {
	switch st := stype.(type) {
	case One:
		return One{ID: core.New[AR]()}
	case TpRef:
		return TpRef{ID: core.New[AR](), Name: st.Name}
	default:
		panic(ErrUnexpectedSt)
	}
}

// Port
type repo[T AR] interface {
	Insert(T) error
	SelectById(core.ID[AR]) (T, error)
	SelectAll() ([]T, error)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend To.*
var (
	ToTpTeaser func(TpRoot) TpTeaser
)

func ToSame(id core.ID[AR]) core.ID[AR] {
	return id
}

func ToCore(id string) (core.ID[AR], error) {
	return core.FromString[AR](id)
}

func ToEdge(id core.ID[AR]) string {
	return id.String()
}
