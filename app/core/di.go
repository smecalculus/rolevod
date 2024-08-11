package core

import (
	"go.uber.org/fx"

	"smecalculus/rolevod/app/core/env"
)

var Module = fx.Module("app/core",
	fx.Provide(
		fx.Annotate(env.NewService, fx.As(new(env.Api))),
	),
)
