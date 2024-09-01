//go:build !goverter

package env

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/env",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(EnvApi))),
	),
	fx.Provide(
		fx.Private,
		newHandlerEcho,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		fx.Annotate(newEnvRepoPgx, fx.As(new(envRepo))),
		fx.Annotate(newTpRepoPgx, fx.As(new(tpRepo))),
	),
	fx.Invoke(
		cfgEcho,
	),
)

//go:embed all:view
var envFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("env").Funcs(sprig.FuncMap()).ParseFS(envFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/api/v1/envs", h.ApiPostOne)
	e.GET("/api/v1/envs/:id", h.ApiGetOne)
	e.GET("/ssr/envs/:id", h.SsrGetOne)
	return nil
}
