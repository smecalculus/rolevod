package env

import (
	"log/slog"

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

func newHandler(a Api, l *slog.Logger) *handlerEcho {
	t := slog.String("t", "env.handlerEcho")
	return &handlerEcho{a, l.With(t)}
}

func newService(r repo, l *slog.Logger) *service {
	t := slog.String("t", "env.service")
	return &service{r, l.With(t)}
}

func newRepo(p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	t := slog.String("t", "env.repoPgx")
	return &repoPgx{p, l.With(t)}
}

func configureEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/envs", h.post)
	e.GET("/envs/:id", h.get)
	return nil
}
