//go:build !goverter

package step

import (
	"go.uber.org/fx"
)

var Module = fx.Module("internal/step",
	fx.Provide(
		fx.Annotate(newRepoPgx[Process], fx.As(new(Repo[Process]))),
		fx.Annotate(newRepoPgx[Message], fx.As(new(Repo[Message]))),
		fx.Annotate(newRepoPgx[Service], fx.As(new(Repo[Service]))),
	),
)
