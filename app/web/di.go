package web

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"

	web "smecalculus/rolevod/app/web/msg"
)

var Module = fx.Module("app/web",
	fx.Provide(
		fx.Private,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		web.NewHandlerEcho,
	),
	fx.Invoke(
		cfgEcho,
	),
)

//go:embed all:msg/feature all:msg/component
var webFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.ParseFS(webFs, "*/*/*.go.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgEcho(e *echo.Echo, h *web.HandlerEcho) {
	e.GET("/", h.Home)
}
