//go:build !goverter

package main

import (
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/msg"

	"smecalculus/rolevod/app/agent"
	"smecalculus/rolevod/app/deal"
	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/seat"
	"smecalculus/rolevod/app/web"
)

func main() {
	fx.New(
		// lib
		core.Module,
		data.Module,
		msg.Module,
		// app
		agent.Module,
		deal.Module,
		role.Module,
		seat.Module,
		web.Module,
	).Run()
}
