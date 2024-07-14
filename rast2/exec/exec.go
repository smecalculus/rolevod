package exec

import (
	"errors"
	"fmt"
	a "smecalculus/rolevod/rast2/ast"
)

type Proc struct {
	C a.Chan
	P a.Expression
}

type Msg struct {
	D a.Chan
	M a.Msg
}

type Sem interface {
	c() a.Chan
}

func (p Proc) c() a.Chan { return p.C }
func (m Msg) c() a.Chan  { return m.D }

type Configuration struct {
	conf   map[a.Chan]Sem
	conts  map[a.Chan]a.Chan
	shared map[a.Chan]a.Chan
}

type Pol string

const (
	Pos Pol = "pos"
	Neg Pol = "neg"
)

var chanNum uint32 = 0

func lfresh() a.Chan {
	n := chanNum
	chanNum += 1
	return a.Chan{fmt.Sprintf("ch%d", n), a.Unknown}
}

func createConfig(sem Sem) *Configuration {
	conf := make(map[a.Chan]Sem)
	conf[sem.c()] = sem
	return &Configuration{conf, make(map[a.Chan]a.Chan), make(map[a.Chan]a.Chan)}
}

func step(env a.Environment, config *Configuration) (*Configuration, error) {
	stepped := true
	for stepped {
		stepped = false
		// итерируем по тому, что меняется на лету
		for _, sem := range config.conf {
			changed, err := matchAndOneStep(env, sem, config)
			if err != nil {
				return nil, err
			}
			if changed {
				stepped = true
			}
		}
	}
	return config, nil
}

func matchAndOneStep(env a.Environment, sem Sem, config *Configuration) (bool, error) {
	switch s := sem.(type) {
	case Proc:
		switch s.P.(type) {
		case a.Fwd:
			return fwd(s.C, config)
		case a.Spawn:
			return spawn(env, s.C, config)
		case a.ExpName:
			return expand(env, s.C, config)
		}
		panic(a.ErrUnexpectedExp)
	case Msg:
		return false, nil
	}
	panic(ErrUnexpectedSem)
}

func fwd(ch a.Chan, config *Configuration) (bool, error) {
	s, ok := config.conf[ch]
	if !ok {
		return false, ErrExecImpossible
	}
	switch proc := s.(type) {
	case Proc:
		fwd, ok := proc.P.(a.Fwd)
		if !ok {
			return false, ErrExecImpossible
		}
		if proc.C != fwd.X {
			return false, ErrChannelMismatch
		}
		// try to apply fwd+ rule
		msgp, err := findMsg(fwd.Y, config, Pos)
		if err != nil {
			return false, err
		}
		if msgp == nil {
			// try to apply fwd- rule
			msgn, err := findMsg(proc.C, config, Neg)
			if err != nil {
				return false, err
			}
			switch msg := msgn.(type) {
			case Msg:
				switch m := msg.M.(type) {
				case a.MLabI, a.MSendL, a.MClose:
					delete(config.conf, proc.C)
					e, ok := config.conts[proc.C]
					if !ok {
						return false, ErrExecImpossible
					}
					delete(config.conf, e)
					msg := Msg{e, a.Msubst(fwd.Y, proc.C, m)}
					config.conf[e] = msg
					delete(config.conts, proc.C)
					addCont(fwd.Y, e, config)
					return true, nil
				}
			default:
				return false, nil
			}
		}
		switch msg := msgp.(type) {
		case Msg:
			switch m := msg.M.(type) {
			case a.MLabI, a.MSendL:
				delete(config.conf, proc.C)
				delete(config.conf, fwd.Y)
				msg := Msg{proc.C, a.Msubst(proc.C, fwd.Y, m)}
				config.conf[proc.C] = msg
				d, ok := config.conts[fwd.Y]
				if !ok {
					return false, ErrExecImpossible
				}
				delete(config.conts, fwd.Y)
				addCont(proc.C, d, config)
				return true, nil
			}
		default:
			return false, nil
		}
	}
	return false, ErrUnexpectedSem
}

func spawn(env a.Environment, ch a.Chan, config *Configuration) (bool, error) {
	s, ok := config.conf[ch]
	if !ok {
		return false, ErrExecImpossible
	}
	delete(config.conf, ch)
	switch sem := s.(type) {
	case Proc:
		exp, ok := sem.P.(a.Spawn)
		if !ok {
			return false, ErrExecImpossible
		}
		c := lfresh()
		proc1 := Proc{c, a.ExpName{c, exp.F, exp.Xs}}
		proc2 := Proc{sem.C, a.Subst(c, exp.X, exp.Q)}
		config.conf[c] = proc1
		config.conf[sem.C] = proc2
		return true, nil
	}
	return false, ErrExecImpossible
}

func expand(env a.Environment, ch a.Chan, config *Configuration) (bool, error) {
	s, ok := config.conf[ch]
	if !ok {
		return false, ErrExecImpossible
	}
	delete(config.conf, ch)
	switch sem := s.(type) {
	case Proc:
		exp, ok := sem.P.(a.ExpName)
		if !ok {
			return false, ErrExecImpossible
		}
		p, err := expdDef(env, exp.X, exp.F, exp.Xs)
		if err != nil {
			return false, err
		}
		proc := Proc{sem.C, a.Subst(sem.C, exp.X, p)}
		config.conf[sem.C] = proc
		return true, nil
	}
	return false, ErrExecImpossible
}

func expdDef(env a.Environment, x a.Chan, f a.Expname, xs []a.Chan) (a.Expression, error) {
	declaration, ok := env[f]
	if !ok {
		return nil, ErrUndefinedProcess
	}
	decdef, ok := declaration.(a.ExpDecDef)
	if !ok {
		return nil, ErrExecImpossible
	}
	exp := a.Subst(x, decdef.ChanTp.Z, decdef.P)
	exp = a.SubstCtx(xs, fst(decdef.Ctx.Ordered), exp)
	return exp, nil
}

func fst(as []a.ChanTp) []a.Chan {
	bs := make([]a.Chan, len(as))
	for _, a := range as {
		bs = append(bs, a.Z)
	}
	return bs
}

func addCont(c1 a.Chan, c2 a.Chan, config *Configuration) {
	cs, ok := config.shared[c1]
	if ok {
		delete(config.shared, c1)
		config.shared[c2] = cs
	}
	config.conts[c1] = c2
}

func findMsg(c1 a.Chan, config *Configuration, dual Pol) (Sem, error) {
	switch dual {
	case Neg:
		c2, ok := config.conts[c1]
		if !ok {
			return nil, nil
		}
		return config.conf[c2], nil
	case Pos:
		return config.conf[c1], nil
	default:
		return nil, ErrExecImpossible
	}
}

func Exec(env a.Environment, f a.Expname) *Configuration {
	c := lfresh()
	sem := Proc{c, a.ExpName{c, f, []a.Chan{}}}
	config, err := step(env, createConfig(sem))
	if err != nil {
		panic(err)
	}
	return config
}

var (
	ErrExecImpossible   = errors.New("exec impossible")
	ErrChannelMismatch  = errors.New("channel mismatch")
	ErrUndefinedProcess = errors.New("udefined process")
	ErrUnexpectedSem    = errors.New("unexpected sem")
)
