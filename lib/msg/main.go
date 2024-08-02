package msg

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module("msg",
	fx.Provide(
		newEcho,
	),
)

func newEcho(lc fx.Lifecycle) *echo.Echo {
	echo := echo.New()
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go echo.Start(":8080")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return echo.Shutdown(ctx)
			},
		},
	)
	return echo
}
