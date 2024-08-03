package msg

import (
	"context"
	"fmt"
	"smecalculus/rolevod/lib/cfg"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module("msg",
	fx.Provide(
		newCfg,
		newEcho,
	),
)

func newCfg(keeper cfg.Keeper) (*props, error) {
	props := &props{}
	err := keeper.Load("msg", props)
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
