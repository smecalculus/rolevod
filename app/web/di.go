package web

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("web",
	fx.Provide(
		fx.Private,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		newHandler,
	),
	fx.Invoke(
		cfgEcho,
	),
)

//go:embed all:feature all:component
var webFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.ParseFS(webFs, "feature/*.go.tmpl", "component/*.go.tmpl")
	if err != nil {
		return nil, err
	}
	name := slog.String("name", "web.rendererStdlib")
	return &msg.RendererStdlib{Registry: t, Log: l.With(name)}, nil
}

func newHandler(r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "web.handlerEcho")
	return &handlerEcho{r, l.With(name)}
}

func cfgEcho(e *echo.Echo, h *handlerEcho) error {
	e.GET("/", h.home)
	return nil
}
