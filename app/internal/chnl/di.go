//go:build !goverter

package chnl

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(Api))),
	),
	fx.Provide(
		fx.Private,
		newHandlerEcho,
		fx.Annotate(newRepoPgx, fx.As(new(repo))),
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
	),
	fx.Invoke(
		cfgEcho,
	),
)

//go:embed all:view
var viesFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("").Funcs(sprig.FuncMap()).ParseFS(viesFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgEcho(e *echo.Echo, h *HandlerEcho) error {
	e.POST("/api/v1/s", h.ApiPostOne)
	e.GET("/api/v1/s/:id", h.ApiGetOne)
	e.GET("/ssr/s/:id", h.SsrGetOne)
	return nil
}
