package ast

type Tp interface {
	tp()
}

func (Plus) tp()   {}
func (With) tp()   {}
func (Tensor) tp() {}
func (Lolli) tp()  {}
func (One) tp()    {}

type Plus struct {
	Choices
}

type With struct {
	Choices
}

type Tensor struct {
	Tp1 Tp
	Tp2 Tp
}

type Lolli struct {
	Tp1 Tp
	Tp2 Tp
}

type One struct{}

type Label string
type Tpname string
type Expname string
type Channel string
type Choices map[Label]Tp

type ChanTp struct {
	Ch Channel
	Tp Tp
}

type Context []ChanTp

type Decl interface {
	decl()
}

func (TpDef) decl()  {}
func (ExpDec) decl() {}
func (Exec) decl()   {}

type TpDef struct {
	Name Tpname
	Tp   Tp
}

type ExpDec struct {
	F           Expname
	Antecedents Context
	Succedent   ChanTp
}

type Exec struct {
	F Expname
}

type Env struct {
	TpDefs  map[Tpname]TpDef
	ExpDecs map[Expname]ExpDec
}

type Exp interface {
	exp()
}

func (Id) exp()      {}
func (Spawn) exp()   {}
func (ExpName) exp() {}
func (Lab) exp()     {}
func (Case) exp()    {}
func (Send) exp()    {}
func (Recv) exp()    {}
func (Close) exp()   {}
func (Wait) exp()    {}
func (Imposs) exp()  {}

type Id struct {
	X Channel
	Y Channel
}

type Spawn struct {
	P Exp
	Q Exp
}

type ExpName struct {
	X  Channel
	F  Expname
	Ys []Channel
}

type Lab struct {
	X Channel
	K Label
	P Exp
}

type Case struct {
	X        Channel
	Branches Branches
}

type Send struct {
	X Channel
	Y Channel
	P Exp
}

type Recv struct {
	X Channel
	Y Channel
	P Exp
}

type Close struct {
	X Channel
}

type Wait struct {
	X Channel
	P Exp
}

type Imposs struct{}

type Branches map[Label]Exp
