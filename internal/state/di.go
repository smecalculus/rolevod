//go:build !goverter

package state

import (
	"go.uber.org/fx"
)

var Module = fx.Module("internal/state",
	fx.Provide(
		fx.Annotate(newRepoPgx, fx.As(new(Repo))),
	),
)
