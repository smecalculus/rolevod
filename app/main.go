//go:build !goverter

package main

import (
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/msg"

	"smecalculus/rolevod/app/force"
	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/seat"
	"smecalculus/rolevod/app/web"
	"smecalculus/rolevod/app/work"
)

func main() {
	fx.New(
		// lib
		core.Module,
		data.Module,
		msg.Module,
		// app
		force.Module,
		role.Module,
		seat.Module,
		web.Module,
		work.Module,
	).Run()
}
