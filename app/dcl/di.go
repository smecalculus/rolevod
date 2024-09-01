//go:build !goverter

package dcl

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/dcl",
	fx.Provide(
		fx.Annotate(newTpService, fx.As(new(TpApi))),
		fx.Annotate(newExpService, fx.As(new(ExpApi))),
	),
	fx.Provide(
		fx.Private,
		newTpHandlerEcho,
		newExpHandlerEcho,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		fx.Annotate(newTpRepoPgx, fx.As(new(repo[TpRoot]))),
		fx.Annotate(newExpRepoPgx, fx.As(new(repo[ExpRoot]))),
	),
	fx.Invoke(
		cfgTpEcho,
		cfgExpEcho,
	),
)

//go:embed all:view
var declFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("dcl").Funcs(sprig.FuncMap()).ParseFS(declFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgTpEcho(e *echo.Echo, h *tpHandlerEcho) error {
	e.POST("/api/v1/tps", h.ApiPostOne)
	e.PUT("/api/v1/tps/:id", h.ApiPutOne)
	e.GET("/ssr/tps/:id", h.SsrGetOne)
	return nil
}

func cfgExpEcho(e *echo.Echo, h *expHandlerEcho) error {
	e.GET("/ssr/exps/:id", h.SsrGetOne)
	return nil
}
