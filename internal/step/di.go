//go:build !goverter

package step

import (
	"go.uber.org/fx"
)

var Module = fx.Module("app/step",
	fx.Provide(
		fx.Annotate(newRepoPgx, fx.As(new(Repo))),
	),
)
