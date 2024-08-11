package main

import (
	"go.uber.org/fx"

	lc "smecalculus/rolevod/lib/core"
	lm "smecalculus/rolevod/lib/msg"
	ls "smecalculus/rolevod/lib/store"

	ac "smecalculus/rolevod/app/core"
	am "smecalculus/rolevod/app/msg"
	as "smecalculus/rolevod/app/store"
	aw "smecalculus/rolevod/app/web"
)

func main() {
	fx.New(
		// lib
		lc.Module,
		lm.Module,
		ls.Module,
		// app
		ac.Module,
		am.Module,
		as.Module,
		aw.Module,
	).Run()
}
