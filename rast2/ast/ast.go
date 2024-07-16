package ast

import (
	"errors"
)

type Mode string

const (
	Shared  Mode = "shared"
	Linear  Mode = "linear"
	Unknown Mode = "unknown"
)

type Label string
type Tpname string
type Expname string
type Choices map[Label]Stype

type Chan struct {
	Name string
	Mode
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
	St1 Stype
	St2 Stype
	Mode
}

type Lolli struct {
	St1 Stype
	St2 Stype
	Mode
}

type One struct{}

type TpName struct {
	Tpname
}

type Up struct {
	St Stype
}

type Down struct {
	St Stype
}

type ChanTp struct {
	Z Chan
	C Stype
}

type Context struct {
	Shared  []ChanTp
	Linear  []ChanTp
	Ordered []ChanTp
}

type Branches map[Label]Expression

type Expression interface {
	expression()
}

func (Fwd) expression()     {}
func (Spawn) expression()   {}
func (ExpName) expression() {}
func (Lab) expression()     {}
func (Case) expression()    {}
func (Send) expression()    {}
func (Recv) expression()    {}
func (Close) expression()   {}
func (Wait) expression()    {}
func (Acquire) expression() {}
func (Accept) expression()  {}
func (Release) expression() {}
func (Detach) expression()  {}

type Fwd struct {
	X Chan
	Y Chan
}

type Spawn struct {
	X  Chan
	F  Expname
	Xs []Chan
	Q  Expression
}

type ExpName struct {
	X  Chan
	F  Expname
	Xs []Chan
}

type Lab struct {
	Chan
	Label
	Expression
}

type Case struct {
	Chan
	Branches
}

type Send struct {
	Ch1 Chan
	Ch2 Chan
	Expression
}

type Recv struct {
	Ch1 Chan
	Ch2 Chan
	Expression
}

type Close struct {
	Chan
}

type Wait struct {
	Chan
	Expression
}

type Acquire struct {
	Ch1 Chan
	Ch2 Chan
	Expression
}

type Accept struct {
	Ch1 Chan
	Ch2 Chan
	Expression
}

type Release struct {
	Ch1 Chan
	Ch2 Chan
	Expression
}

type Detach struct {
	Ch1 Chan
	Ch2 Chan
	Expression
}

type Decl interface {
	decl()
}

func (TpDef) decl()     {}
func (ExpDecDef) decl() {}
func (Exec) decl()      {}

type TpDef struct {
	A  Tpname
	St Stype
}

type ExpDecDef struct {
	F      Expname
	M      Mode
	Ctx    Context
	ChanTp ChanTp
	P      Expression
}

type Exec struct {
	F Expname
}

type Environment struct {
	TpDefs     map[Tpname]TpDef
	ExpDecDefs map[Expname]ExpDecDef
}

type Msg interface {
	msg()
}

func (MLabI) msg()  {}
func (MLabE) msg()  {}
func (MSendT) msg() {}
func (MSendL) msg() {}
func (MClose) msg() {}

type MLabI struct {
	X Chan
	K Label
	Y Chan
}

type MLabE struct {
	X Chan
	K Label
	Y Chan
}

type MSendT struct {
	X Chan
	W Chan
	Y Chan
}

type MSendL struct {
	X Chan
	W Chan
	Y Chan
}

type MClose struct {
	X Chan
}

func sub(new Chan, old Chan, x Chan) Chan {
	if x.Name == old.Name {
		return new
	}
	return x
}

func Subst(new Chan, old Chan, expr Expression) Expression {
	switch exp := expr.(type) {
	case Fwd:
		return Fwd{sub(new, old, exp.X), sub(new, old, exp.Y)}
	case Spawn:
		return Spawn{exp.X, exp.F, SubstList(new, old, exp.Xs), Subst(new, old, exp.Q)}
	case ExpName:
		return ExpName{exp.X, exp.F, SubstList(new, old, exp.Xs)}
	}
	panic(ErrUnexpectedExp)
}

func Msubst(new Chan, old Chan, m Msg) Msg {
	switch msg := m.(type) {
	case MLabI:
		return MLabI{sub(new, old, msg.X), msg.K, sub(new, old, msg.Y)}
	case MLabE:
		return MLabI{sub(new, old, msg.X), msg.K, sub(new, old, msg.Y)}
	case MSendT:
		return MSendT{sub(new, old, msg.X), msg.W, sub(new, old, msg.Y)}
	case MSendL:
		return MSendT{sub(new, old, msg.X), msg.W, sub(new, old, msg.Y)}
	case MClose:
		return MClose{sub(new, old, msg.X)}
	}
	panic(ErrUnexpectedMsg)
}

func SubstList(new Chan, old Chan, xs []Chan) []Chan {
	ys := make([]Chan, len(xs))
	for _, x := range xs {
		ys = append(ys, sub(new, old, x))
	}
	return ys
}

func SubstCtx(ctx1 []Chan, ctx2 []Chan, expr Expression) Expression {
	if len(ctx1) != len(ctx2) {
		panic(ErrAstImpossible)
	}
	exp := expr
	for i, c1 := range ctx1 {
		c2 := ctx2[i]
		exp = Subst(c1, c2, exp)
	}
	return exp
}

var (
	ErrUnexpectedExp = errors.New("unexpected exp")
	ErrUnexpectedMsg = errors.New("unexpected msg")
	ErrAstImpossible = errors.New("ast impossible")
)
