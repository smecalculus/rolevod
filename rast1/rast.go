package rast

import (
	"errors"
)

//

type tp interface {
	tp()
}

func (Plus) tp()   {}
func (With) tp()   {}
func (Tensor) tp() {}
func (Lolli) tp()  {}
func (One) tp()    {}

type Plus map[label]tp

type With map[label]tp

type Tensor struct {
	tp1 tp
	tp2 tp
}

type Lolli struct {
	tp1 tp
	tp2 tp
}

type One struct{}

//

type label string
type tpname string
type expname string
type channel string

type choices map[label]tp

type chan_tp struct {
	ch channel
	tp tp
}

type context []chan_tp

//

type decl interface {
	decl()
}

func (TpDef) decl()  {}
func (ExpDec) decl() {}

type TpDef struct {
	name tpname
	tp   tp
}

type ExpDec struct {
	name        expname
	antecedents context
	succedent   chan_tp
}

//

type env struct {
	tpDefs  map[tpname]TpDef
	expDecs map[expname]ExpDec
}

//

type exp interface {
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
	X channel
	Y channel
}

type Spawn struct {
	P exp
	Q exp
}

type ExpName struct {
	X  channel
	F  expname
	Ys []channel
}

type Lab struct {
	X channel
	K label
	P exp
}

type Case struct {
	X        channel
	Branches branches
}

type Send struct {
	X channel
	Y channel
	P exp
}

type Recv struct {
	X channel
	Y channel
	P exp
}

type Close struct {
	X channel
}

type Wait struct {
	X channel
	P exp
}

type Imposs struct{}

//

type branches map[label]exp

//

type value interface {
	value()
}

func (LabV) value()    {}
func (SendV) value()   {}
func (CloseV) value()  {}
func (CloCase) value() {}
func (CloRecv) value() {}

type LabV struct {
	K label
	V value
}

type SendV struct {
	W value
	V value
}

type CloseV struct{}

type CloCase struct {
	Eta      eta
	Branches branches
	Z        channel
}

type chan_exp struct {
	X channel
	P exp
}

type CloRecv struct {
	Eta     eta
	ChanExp chan_exp
	Z       channel
}

//

type eta map[channel]value

func evaluate(env env, P exp, z channel) (value, error) {
	return eval(env, make(eta), P, z)
}

func eval(env env, eta eta, P exp, z channel) (value, error) {
	switch exp := P.(type) {
	case Id:
		return eta[exp.Y], nil
	case Spawn:
		eta2, err := eval_call(env, eta, exp.P)
		if err != nil {
			return nil, err
		}
		return eval(env, eta2, exp.Q, z)
	case Lab:
		if exp.X == z {
			v, err := eval(env, eta, exp.P, z)
			if err != nil {
				return nil, err
			}
			return LabV{exp.K, v}, nil
		}
		clo := eta[exp.X].(CloCase)
		v, err := eval(env, clo.Eta, clo.Branches[exp.K], z)
		if err != nil {
			return nil, err
		}
		eta[exp.X] = v
		return eval(env, eta, exp.P, z)
	case Case:
		if exp.X == z {
			return CloCase{eta, exp.Branches, z}, nil
		}
		lab := eta[exp.X].(LabV)
		eta[exp.X] = lab.V
		return eval(env, eta, exp.Branches[lab.K], z)
	case Send:
		w := eta[exp.Y]
		if exp.X == z {
			delete(eta, exp.Y)
			v, err := eval(env, eta, exp.P, z)
			if err != nil {
				return nil, err
			}
			return SendV{w, v}, nil
		}
		clo := eta[exp.X].(CloRecv)
		clo.Eta[clo.ChanExp.X] = w
		v, err := eval(env, clo.Eta, clo.ChanExp.P, z)
		if err != nil {
			return nil, err
		}
		eta[exp.X] = v
		delete(eta, exp.Y)
		return eval(env, eta, exp.P, z)
	case Recv:
		if exp.X == z {
			return CloRecv{eta, chan_exp{exp.Y, exp.P}, z}, nil
		}
		send := eta[exp.X].(SendV)
		eta[exp.X] = send.V
		eta[exp.Y] = send.W
		return eval(env, eta, exp.P, z)
	case Close:
		return CloseV{}, nil
	case Wait:
		_, ok := eta[exp.X].(CloseV)
		if !ok {
			return nil, ErrUnexpectedWait
		}
		delete(eta, exp.X)
		return eval(env, eta, exp.P, z)
	case nil:
		return nil, nil
	default:
		return nil, ErrUnexpectedExp
	}
}

var (
	ErrUnexpectedExp  = errors.New("unexpected exp")
	ErrUnexpectedWait = errors.New("unexpected wait")
)

// Simplified version which based on ExpDec instead of ExpDef
func eval_call(env env, eta1 eta, P exp) (eta, error) {
	switch exp := P.(type) {
	case ExpName:
		expDec := env.expDecs[exp.F]
		ys := make([]channel, len(expDec.antecedents))
		for i, ant := range expDec.antecedents {
			ys[i] = ant.ch
		}
		if len(exp.Ys) != len(ys) {
			return nil, errors.New("lenghts should be equal")
		}
		eta2 := make(eta, len(ys))
		for i, y := range exp.Ys {
			eta2[ys[i]] = eta1[y]
			delete(eta1, y)
		}
		v, err := eval(env, eta1, P, expDec.succedent.ch)
		if err != nil {
			return nil, err
		}
		eta2[exp.X] = v
		return eta2, nil
	default:
		return nil, ErrUnexpectedExp
	}
}

type Rastovod interface {
	eval() (error)
}

type rastovodImpl struct {
	env env
	eta eta
}