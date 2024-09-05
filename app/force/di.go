//go:build !goverter

package force

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/force",
	fx.Provide(
		fx.Annotate(newForceService, fx.As(new(ForceApi))),
	),
	fx.Provide(
		fx.Private,
		newForceHandlerEcho,
		fx.Annotate(newForceRepoPgx, fx.As(new(forceRepo))),
		newKinshipHandlerEcho,
		fx.Annotate(newKinshipRepoPgx, fx.As(new(kinshipRepo))),
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
	),
	fx.Invoke(
		cfgForceEcho,
		cfgKinshipEcho,
	),
)

//go:embed all:view
var viesFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("force").Funcs(sprig.FuncMap()).ParseFS(viesFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgForceEcho(e *echo.Echo, h *forceHandlerEcho) error {
	e.POST("/api/v1/forces", h.ApiPostOne)
	e.GET("/api/v1/forces/:id", h.ApiGetOne)
	e.GET("/ssr/forces/:id", h.SsrGetOne)
	return nil
}

func cfgKinshipEcho(e *echo.Echo, h *kinshipHandlerEcho) error {
	e.POST("/api/v1/forces/:id/kinships", h.ApiPostOne)
	return nil
}
