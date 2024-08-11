package msg

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"

	me "smecalculus/rolevod/app/msg/env"
)

var Module = fx.Module("app/msg",
	fx.Provide(
		// fx.Private,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		me.NewHandlerEcho,
	),
	fx.Invoke(
		cfgEcho,
	),
)

//go:embed env/*.go.html
var msgFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.ParseFS(msgFs, "env/*.go.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgEcho(e *echo.Echo, h *me.HandlerEcho) error {
	e.POST("/envs", h.Post)
	e.GET("/envs/:id", h.Get)
	return nil
}
