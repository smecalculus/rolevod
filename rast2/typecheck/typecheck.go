package typecheck

import (
	"fmt"
	a "smecalculus/rolevod/rast2/ast"
)

func CheckExp(env a.Environment, delta a.Context, exp a.Expression, zc a.ChanTp) error {
	return fmt.Errorf("not implemented yet")
}

func Contractive(a a.Stype) bool {
	return false
}

func EsyncTp(env a.Environment, tp a.Stype) error {
	return fmt.Errorf("not implemented yet")
}
