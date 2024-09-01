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

// Aggregate Spec
type AS interface {
	as()
}

func (TpSpec) as()  {}
func (ExpSpec) as() {}

type TpSpec struct {
	Name Tpname
}

type ExpSpec struct {
	Name Expname
}

// Aggregate Root (aka decl)
type AR interface {
	ar()
}

func (TpRoot) ar()  {}
func (ExpRoot) ar() {}

// aka TpDef
type TpRoot struct {
	ID   core.ID[AR]
	Name Tpname
	St   Stype
}

type TpTeaser struct {
	ID   core.ID[AR]
	Name Tpname
}

// aka ExpDecDef
type ExpRoot struct {
	ID   core.ID[AR]
	Name Expname
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
	A Stype
}

type Down struct {
	A Stype
}

type ChanTp struct {
	X Chan
	A Stype
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
	}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *tpService) Update(root TpRoot) error {
	err := s.repo.Insert(root)
	if err != nil {
		return err
	}
	return nil
}

func (s *tpService) Retrieve(id core.ID[AR]) (TpRoot, error) {
	root, err := s.repo.SelectById(id)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *tpService) RetreiveAll() ([]TpRoot, error) {
	roots, err := s.repo.SelectAll()
	if err != nil {
		return nil, err
	}
	return roots, nil
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
	root, err := s.repo.SelectById(id)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *expService) RetreiveAll() ([]ExpRoot, error) {
	roots, err := s.repo.SelectAll()
	if err != nil {
		return nil, err
	}
	return roots, nil
}

// Port
type repo[T AR] interface {
	Insert(T) error
	SelectById(core.ID[AR]) (T, error)
	SelectAll() ([]T, error)
}

func toCore(id string) (core.ID[AR], error) {
	return core.FromString[AR](id)
}

func toEdge(id core.ID[AR]) string {
	return core.ToString(id)
}
