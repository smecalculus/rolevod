package env

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"

	ec "smecalculus/rolevod/app/env/core"
	em "smecalculus/rolevod/app/env/msg"
	es "smecalculus/rolevod/app/env/store"
)

var Module = fx.Module("app/env",
	fx.Provide(
		fx.Annotate(ec.NewService, fx.As(new(ec.Api))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		em.NewHandlerEcho,
		fx.Annotate(es.NewRepoPgx, fx.As(new(ec.Repo))),
	),
	fx.Invoke(
		cfgEcho,
	),
)

//go:embed msg/*.go.html
var envFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.ParseFS(envFs, "msg/*.go.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgEcho(e *echo.Echo, h *em.HandlerEcho) error {
	e.POST("/envs", h.Post)
	e.GET("/envs/:id", h.Get)
	return nil
}
