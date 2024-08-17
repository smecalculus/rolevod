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
		fx.Annotate(newService, fx.As(new(Api))),
	),
	fx.Provide(
		fx.Private,
		newHandlerEcho,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		fx.Annotate(newRepoPgx, fx.As(new(repo))),
	),
	fx.Invoke(
		cfgEcho,
	),
)

//go:embed all:view
var declFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("decl").Funcs(sprig.FuncMap()).ParseFS(declFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgEcho(e *echo.Echo, h *handlerEcho) error {
	e.GET("/ssr/decls/:id", h.SsrGetOne)
	return nil
}
