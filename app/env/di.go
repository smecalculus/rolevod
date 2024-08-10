package env

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("env",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(Api))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		newHandler,
		fx.Annotate(newRepo, fx.As(new(repo))),
	),
	fx.Invoke(
		configureEcho,
	),
)

//go:embed *.go.tmpl
var envFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.ParseFS(envFs, "*.go.tmpl")
	if err != nil {
		return nil, err
	}
	name := slog.String("name", "env.rendererStdlib")
	return &msg.RendererStdlib{Registry: t, Log: l.With(name)}, nil
}

func newHandler(a Api, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "env.handlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func newService(r repo, l *slog.Logger) *service {
	name := slog.String("name", "env.service")
	return &service{r, l.With(name)}
}

func newRepo(p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	name := slog.String("name", "env.repoPgx")
	return &repoPgx{p, l.With(name)}
}

func configureEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/envs", h.post)
	e.GET("/envs/:id", h.get)
	return nil
}
