package msg

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
)

var Module = fx.Module("msg",
	fx.Provide(
		newEcho,
	),
	fx.Provide(
		fx.Private,
		newCfg,
	),
)

func newCfg(k core.Keeper) (*props, error) {
	props := &props{}
	err := k.Load("msg", props)
	if err != nil {
		return nil, err
	}
	return props, nil
}

func newEcho(props *props, lc fx.Lifecycle) *echo.Echo {
	echo := echo.New()
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go echo.Start(fmt.Sprintf(":%v", props.Protocol.Http.Port))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return echo.Shutdown(ctx)
			},
		},
	)
	return echo
}
