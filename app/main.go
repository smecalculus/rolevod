//go:build !goverter

package main

import (
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
	"smecalculus/rolevod/lib/store"

	"smecalculus/rolevod/app/decl"
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
		web.Module,
		env.Module,
		decl.Module,
	).Run()
}
