//go:build !goverter

package step

import (
	"go.uber.org/fx"
)

var Module = fx.Module("internal/step",
	fx.Provide(
		fx.Annotate(newRepoPgx2, fx.As(new(Repo2))),
		fx.Annotate(newRepoPgx[ProcRoot], fx.As(new(Repo[ProcRoot]))),
		fx.Annotate(newRepoPgx[MsgRoot], fx.As(new(Repo[MsgRoot]))),
		fx.Annotate(newRepoPgx[SrvRoot], fx.As(new(Repo[SrvRoot]))),
	),
)
