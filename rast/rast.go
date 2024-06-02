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
	tpname tpname
	tp     tp
}

type ExpDec struct {
	expname     expname
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
	x channel
	y channel
}

type Spawn struct {
	P exp
	Q exp
}

type ExpName struct {
	x  channel
	f  expname
	ys []channel
}

type Lab struct {
	x channel
	k label
	P exp
}

type Case struct {
	x        channel
	branches branches
}

type Send struct {
	x channel
	y channel
	P exp
}

type Recv struct {
	x channel
	y channel
	P exp
}

type Close struct {
	x channel
}

type Wait struct {
	x channel
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
	k label
	v value
}

type SendV struct {
	w value
	v value
}

type CloseV struct{}

type CloCase struct {
	eta      eta
	branches branches
	z        channel
}

type chan_exp struct {
	x channel
	P exp
}

type CloRecv struct {
	eta      eta
	chan_exp chan_exp
	z        channel
}

//

type eta map[channel]value

func evaluate(env env, P exp, z channel) (value, error) {
	return eval(env, make(eta), P, z)
}

func eval(env env, eta eta, P exp, z channel) (value, error) {
	switch exp := P.(type) {
	case Id:
		return eta[exp.y], nil
	case Spawn:
		eta2, err := eval_call(env, eta, exp.P)
		if err != nil {
			return nil, err
		}
		return eval(env, eta2, exp.Q, z)
	case Lab:
		if exp.x == z {
			v, err := eval(env, eta, P, z)
			if err != nil {
				return nil, err
			}
			return LabV{exp.k, v}, nil
		}
		clo := eta[exp.x].(CloCase)
		v, err := eval(env, clo.eta, clo.branches[exp.k], z)
		if err != nil {
			return nil, err
		}
		eta[exp.x] = v
		return eval(env, eta, P, z)
	case Case:
		if exp.x == z {
			return CloCase{eta, exp.branches, z}, nil
		}
		lab := eta[exp.x].(LabV)
		eta[exp.x] = lab.v
		return eval(env, eta, exp.branches[lab.k], z)
	case Send:
		if exp.x == z {
			delete(eta, exp.y)
			v, err := eval(env, eta, P, z)
			if err != nil {
				return nil, err
			}
			// получаем то, что удалено
			return SendV{eta[exp.y], v}, nil
		}
		clo := eta[exp.x].(CloRecv)
		clo.eta[clo.chan_exp.x] = eta[exp.y]
		v, err := eval(env, clo.eta, clo.chan_exp.P, z)
		if err != nil {
			return nil, err
		}
		eta[exp.x] = v
		delete(eta, exp.y)
		return eval(env, eta, P, z)
	case Recv:
		if exp.x == z {
			return CloRecv{eta, chan_exp{exp.y, exp.P}, z}, nil
		}
		send := eta[exp.x].(SendV)
		eta[exp.x] = send.v
		eta[exp.y] = send.w
		return eval(env, eta, exp.P, z)
	case Close:
		return CloseV{}, nil
	case Wait:
		_, ok := eta[exp.x].(CloseV)
		if !ok {
			return nil, ErrUnexpectedWait
		}
		delete(eta, exp.x)
		return eval(env, eta, exp.P, z)
	case nil:
		return nil, nil
	default:
		return nil, errors.New("unexpected exp")
	}
}

var ErrUnexpectedWait = errors.New("unexpected wait")

// simplified version which based on ExpDec instead of ExpDef
func eval_call(env env, eta1 eta, P exp) (eta, error) {
	switch exp := P.(type) {
	case ExpName:
		expDec := env.expDecs[exp.f]
		ys := make([]channel, len(expDec.antecedents))
		for i, ant := range expDec.antecedents {
			ys[i] = ant.ch
		}
		if len(exp.ys) != len(ys) {
			return nil, errors.New("lenghts should be equal")
		}
		eta2 := make(eta, len(ys))
		for i, y := range exp.ys {
			eta2[ys[i]] = eta1[y]
			delete(eta1, y)
		}
		v, err := eval(env, eta1, P, expDec.succedent.ch)
		if err != nil {
			return nil, err
		}
		eta2[exp.x] = v
		return eta2, nil
	default:
		return nil, errors.New("unexpected exp")
	}
}
