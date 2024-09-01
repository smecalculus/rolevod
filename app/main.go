//go:build !goverter

package main

import (
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/msg"

	"smecalculus/rolevod/app/dcl"
	"smecalculus/rolevod/app/web"
	ws "smecalculus/rolevod/app/ws"
)

func main() {
	fx.New(
		// lib
		core.Module,
		data.Module,
		msg.Module,
		// app
		web.Module,
		ws.Module,
		dcl.Module,
	).Run()
}
