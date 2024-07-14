package eval

import (
	"errors"
	a "smecalculus/rolevod/rast1/ast"
)

type Value interface {
	value()
}

func (Lab) value()     {}
func (Send) value()    {}
func (Close) value()   {}
func (CloCase) value() {}
func (CloRecv) value() {}

type Lab struct {
	K a.Label
	V Value
}

type Send struct {
	W Value
	V Value
}

type Close struct{}

type CloCase struct {
	Eta      Eta
	Branches a.Branches
	Z        a.Channel
}

type CloRecv struct {
	Eta     Eta
	ChanExp ChanExp
	Z       a.Channel
}

type ChanExp struct {
	X a.Channel
	P a.Exp
}

type Eta map[a.Channel]Value

func Evaluate(env a.Env, P a.Exp, z a.Channel) (Value, error) {
	return Eval(env, make(Eta), P, z)
}

func Eval(env a.Env, eta Eta, P a.Exp, z a.Channel) (Value, error) {
	switch exp := P.(type) {
	case a.Id:
		return eta[exp.Y], nil
	case a.Spawn:
		eta2, err := evalCall(env, eta, exp.P)
		if err != nil {
			return nil, err
		}
		return Eval(env, eta2, exp.Q, z)
	case a.Lab:
		if exp.X == z {
			v, err := Eval(env, eta, exp.P, z)
			if err != nil {
				return nil, err
			}
			return Lab{exp.K, v}, nil
		}
		clo := eta[exp.X].(CloCase)
		v, err := Eval(env, clo.Eta, clo.Branches[exp.K], z)
		if err != nil {
			return nil, err
		}
		eta[exp.X] = v
		return Eval(env, eta, exp.P, z)
	case a.Case:
		if exp.X == z {
			return CloCase{eta, exp.Branches, z}, nil
		}
		lab := eta[exp.X].(Lab)
		eta[exp.X] = lab.V
		return Eval(env, eta, exp.Branches[lab.K], z)
	case a.Send:
		w := eta[exp.Y]
		if exp.X == z {
			delete(eta, exp.Y)
			v, err := Eval(env, eta, exp.P, z)
			if err != nil {
				return nil, err
			}
			return Send{w, v}, nil
		}
		clo := eta[exp.X].(CloRecv)
		clo.Eta[clo.ChanExp.X] = w
		v, err := Eval(env, clo.Eta, clo.ChanExp.P, z)
		if err != nil {
			return nil, err
		}
		eta[exp.X] = v
		delete(eta, exp.Y)
		return Eval(env, eta, exp.P, z)
	case a.Recv:
		if exp.X == z {
			return CloRecv{eta, ChanExp{exp.Y, exp.P}, z}, nil
		}
		send := eta[exp.X].(Send)
		eta[exp.X] = send.V
		eta[exp.Y] = send.W
		return Eval(env, eta, exp.P, z)
	case a.Close:
		return Close{}, nil
	case a.Wait:
		_, ok := eta[exp.X].(Close)
		if !ok {
			return nil, ErrUnexpectedWait
		}
		delete(eta, exp.X)
		return Eval(env, eta, exp.P, z)
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
func evalCall(env a.Env, eta1 Eta, P a.Exp) (Eta, error) {
	switch exp := P.(type) {
	case a.ExpName:
		expDec := env.ExpDecs[exp.F]
		ys := make([]a.Channel, len(expDec.Antecedents))
		for i, ant := range expDec.Antecedents {
			ys[i] = ant.Ch
		}
		if len(exp.Ys) != len(ys) {
			return nil, errors.New("lenghts should be equal")
		}
		eta2 := make(Eta, len(ys))
		for i, y := range exp.Ys {
			eta2[ys[i]] = eta1[y]
			delete(eta1, y)
		}
		v, err := Eval(env, eta1, P, expDec.Succedent.Ch)
		if err != nil {
			return nil, err
		}
		eta2[exp.X] = v
		return eta2, nil
	default:
		return nil, ErrUnexpectedExp
	}
}
