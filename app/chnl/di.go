//go:build !goverter

package chnl

import (
	"go.uber.org/fx"
)

var Module = fx.Module("internal/chnl",
	fx.Provide(
		fx.Annotate(newRepoPgx, fx.As(new(Repo))),
	),
)
