package main

import (
	"smecalculus/rolevod/app/env"
	"smecalculus/rolevod/lib/cfg"
	"smecalculus/rolevod/lib/db"
	"smecalculus/rolevod/lib/msg"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		// lib
		cfg.Module,
		msg.Module,
		db.Module,
		// app
		env.Module,
	).Run()
}
