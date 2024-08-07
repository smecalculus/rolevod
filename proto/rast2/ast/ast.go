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
type Tpname = string
type Expname = string
type Choices map[Label]Stype

type Chan struct {
	V string
	M Mode
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
	X Chan
}

type Wait struct {
	X Chan
	P Expression
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
	V Tpname
	A Stype
}

type ExpDecDef struct {
	F   Expname
	Ctx Context
	Zc  ChanTp
	P   Expression
}

type Exec struct {
	F Expname
}

type Environment map[string]Decl

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
	if x.V == old.V {
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
	case Close:
		return Close{sub(new, old, exp.X)}
	case Wait:
		return Wait{sub(new, old, exp.X), Subst(new, old, exp.P)}
	default:
		panic(ErrUnexpectedExp)
	}
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
	default:
		panic(ErrUnexpectedMsg)
	}
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

func IsShared(env Environment, c Stype) bool {
	return false
}

func ExpdTp(env Environment, v string) (Stype, error) {
	decl, ok := env[v].(TpDef)
	if !ok {
		return nil, ErrAstImpossible
	}
	return decl.A, nil
}

var (
	ErrUnexpectedExp = errors.New("unexpected exp")
	ErrUnexpectedMsg = errors.New("unexpected msg")
	ErrAstImpossible = errors.New("ast impossible")
)
