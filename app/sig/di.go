//go:build !goverter

package sig

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/sig",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
		fx.Annotate(newRepoPgx, fx.As(new(Repo))),
	),
	fx.Provide(
		fx.Private,
		newHandlerEcho,
		newKinshipHandlerEcho,
		fx.Annotate(newKinshipRepoPgx, fx.As(new(kinshipRepo))),
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
	),
	fx.Invoke(
		cfgSigEcho,
		cfgKinshipEcho,
	),
)

//go:embed all:view
var viesFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("sig").Funcs(sprig.FuncMap()).ParseFS(viesFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgSigEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/api/v1/sigs", h.ApiPostOne)
	e.GET("/api/v1/sigs/:id", h.ApiGetOne)
	e.GET("/ssr/sigs/:id", h.SsrGetOne)
	return nil
}

func cfgKinshipEcho(e *echo.Echo, h *kinshipHandlerEcho) error {
	e.POST("/api/v1/sigs/:id/kinships", h.ApiPostOne)
	return nil
}
