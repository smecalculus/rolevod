package dcl

import (
	"log/slog"

	"smecalculus/rolevod/lib/core"
)

type Spec struct {
	Name string
}

type Label string
type Choices map[Label]Stype

type Chan struct {
	V string
}

type Stype interface {
	stype()
}

func (Plus) stype()   {}
func (With) stype()   {}
func (Tensor) stype() {}
func (Lolli) stype()  {}
func (One) stype()    {}
func (TpName) stype() {}
func (Up) stype()     {}
func (Down) stype()   {}

type Plus struct {
	Choices
}

type With struct {
	Choices
}

type Tensor struct {
	S Stype
	T Stype
}

type Lolli struct {
	S Stype
	T Stype
}

type One struct{}

type TpName struct {
	A Tpname
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

type Dcl core.Entity

type Root interface {
	decl()
}

func (TpDef) decl()     {}
func (ExpDecDef) decl() {}

type Tpname = string

type TpDef struct {
	ID   core.ID[Dcl]
	Name Tpname
}

type Expname = string

type ExpDecDef struct {
	ID   core.ID[Dcl]
	Name Expname
	Zc   ChanTp
}

// port
type Api interface {
	Create(Spec) (Root, error)
	Retrieve(core.ID[Dcl]) (Root, error)
	RetreiveAll() ([]Root, error)
}

// core
type service struct {
	repo repo
	log  *slog.Logger
}

func newService(r repo, l *slog.Logger) *service {
	name := slog.String("name", "decl.service")
	return &service{r, l.With(name)}
}

func (s *service) Create(spec Spec) (Root, error) {
	root := TpDef{
		ID:   core.New[Dcl](),
		Name: spec.Name,
	}
	err := s.repo.Insert(root)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) Retrieve(id core.ID[Dcl]) (Root, error) {
	root, err := s.repo.SelectById(id)
	if err != nil {
		return root, err
	}
	return root, nil
}

func (s *service) RetreiveAll() ([]Root, error) {
	tpDefs, err := s.repo.SelectAll()
	if err != nil {
		return nil, err
	}
	roots := make([]Root, len(tpDefs))
	for i, v := range tpDefs {
		roots[i] = Root(v)
	}
	return roots, nil
}

// port
type repo interface {
	Insert(TpDef) error
	SelectById(core.ID[Dcl]) (TpDef, error)
	SelectAll() ([]TpDef, error)
}

func ToCore(id string) (core.ID[Dcl], error) {
	return core.FromString[Dcl](id)
}

func ToEdge(id core.ID[Dcl]) string {
	return core.ToString(id)
}

func toCore(id string) (core.ID[Dcl], error) {
	return core.FromString[Dcl](id)
}

func toEdge(id core.ID[Dcl]) string {
	return core.ToString(id)
}
