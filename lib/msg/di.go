package msg

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module("msg",
	fx.Provide(
		newConf,
		newEcho,
	),
)

func newConf() props {
	return props{
		Protocol: protocol{
			Http: http{
				Port: 8080,
			},
		},
	}
}

func newEcho(props props, lc fx.Lifecycle) *echo.Echo {
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
