package elab

import (
	"errors"
	"fmt"
	"slices"
	a "smecalculus/rolevod/rast2/ast"
	tc "smecalculus/rolevod/rast2/typecheck"
)

func ElabTps(env a.Environment, dcls map[string]a.Decl) error {
	for _, v := range dcls {
		switch dcl := v.(type) {
		case a.TpDef:
			if !tc.Contractive(dcl.St) {
				return ErrTypeNotContractive(dcl.St)
			}
			err := tc.EsyncTp(env, dcl.St)
			if err != nil {
				return err
			}
		case a.ExpDecDef:
			delta := slices.Clone(dcl.Ctx.Linear)
			if dups(append(delta, dcl.Zc)) {
				return ErrDuplicateVariable
			}
			err := tc.EsyncTp(env, dcl.Zc.A)
			if err != nil {
				return err
			}
		default:
			continue
		}
	}
	return nil
}

func ElabExps(env a.Environment, dcls map[string]a.Decl) error {
	for _, v := range dcls {
		switch dcl := v.(type) {
		case a.TpDef:
			// already checked validity during first pass
			continue
		case a.ExpDecDef:
			err := tc.CheckExp(env, dcl.Ctx, dcl.P, dcl.Zc)
			if err != nil {
				return err
			}
		case a.Exec:
			dec, ok := env[dcl.F].(a.ExpDecDef)
			if !ok {
				return ErrProcessUndefined(dcl.F)
			}
			if len(dec.Ctx.Ordered) > 0 {
				return ErrProcessNonEmptyContext(dcl.F)
			}
		default:
			return ErrElabImpossible
		}
	}
	return nil
}

func ElabDecls(env a.Environment, dcls map[string]a.Decl) error {
	// first pass: check validity of types and create internal names
	err1 := ElabTps(env, dcls)
	if err1 != nil {
		return err1
	}
	// second pass: perform reconstruction and type checking
	err2 := ElabExps(env, env)
	if err2 != nil {
		return err2
	}
	return nil
}

func dups(delta []a.ChanTp) bool {
	var xs []string
	for _, v := range delta {
		xs = append(xs, v.X.Id)
	}
	slices.Sort(xs)
	for i, v := range xs {
		if slices.Contains(xs[i:], v) {
			return true
		}
	}
	return false
}

var (
	ErrElabImpossible    = errors.New("elab impossible")
	ErrDuplicateVariable = errors.New("duplicate variable in process declaration")
)

func ErrProcessUndefined(f string) error {
	return fmt.Errorf("process %q udefined", f)
}

func ErrProcessNonEmptyContext(f string) error {
	return fmt.Errorf("process %q has a non-empty context, cannot be executed", f)
}

func ErrTypeNotContractive(a a.Stype) error {
	return fmt.Errorf("type not contractive: %v", a)
}
