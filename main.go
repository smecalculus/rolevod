package main

import (
	"smecalculus/rolevod/app/env"
	"smecalculus/rolevod/lib/db"
	"smecalculus/rolevod/lib/msg"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		env.Module,
		msg.Module,
		db.Module,
	).Run()
}
