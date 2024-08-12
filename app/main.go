//go:build !goverter

package main

import (
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
	"smecalculus/rolevod/lib/store"

	"smecalculus/rolevod/app/env"
	"smecalculus/rolevod/app/web"
)

func main() {
	fx.New(
		// lib
		core.Module,
		msg.Module,
		store.Module,
		// app
		env.Module,
		web.Module,
	).Run()
}
