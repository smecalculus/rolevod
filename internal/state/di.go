//go:build !goverter

package state

import (
	"go.uber.org/fx"
)

var Module = fx.Module("app/state",
	fx.Provide(
		fx.Annotate(newRepoPgx, fx.As(new(Repo))),
	),
)
