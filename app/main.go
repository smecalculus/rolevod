package main

import (
	"go.uber.org/fx"

	"smecalculus/rolevod/app/env"
	"smecalculus/rolevod/app/web"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/db"
	"smecalculus/rolevod/lib/msg"
)

func main() {
	fx.New(
		// lib
		core.Module,
		msg.Module,
		db.Module,
		// app
		web.Module,
		env.Module,
	).Run()
}
