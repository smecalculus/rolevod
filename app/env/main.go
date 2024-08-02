package env

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module("env",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(Api))),
	),
	fx.Provide(
		fx.Private,
		newHandler,
		fx.Annotate(newRepo, fx.As(new(repo))),
	),
	fx.Invoke(
		configureEcho,
	),
)

func newHandler(api Api) *handlerEcho {
	return &handlerEcho{api}
}

func newService(repo repo) *service {
	return &service{repo}
}

func newRepo(pgx *pgxpool.Pool) *repoPgx {
	return &repoPgx{pgx}
}

func configureEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/envs", h.post)
	e.GET("/envs/:id", h.get)
	return nil
}
