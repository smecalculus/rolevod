//go:build !goverter

package work

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/work",
	fx.Provide(
		fx.Annotate(newWorkService, fx.As(new(WorkApi))),
	),
	fx.Provide(
		fx.Private,
		newWorkHandlerEcho,
		fx.Annotate(newWorkRepoPgx, fx.As(new(workRepo))),
		newKinshipHandlerEcho,
		fx.Annotate(newKinshipRepoPgx, fx.As(new(kinshipRepo))),
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
	),
	fx.Invoke(
		cfgWorkEcho,
		cfgKinshipEcho,
	),
)

//go:embed all:view
var viesFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("work").Funcs(sprig.FuncMap()).ParseFS(viesFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgWorkEcho(e *echo.Echo, h *workHandlerEcho) error {
	e.POST("/api/v1/works", h.ApiPostOne)
	e.GET("/api/v1/works/:id", h.ApiGetOne)
	e.GET("/ssr/works/:id", h.SsrGetOne)
	return nil
}

func cfgKinshipEcho(e *echo.Echo, h *kinshipHandlerEcho) error {
	e.POST("/api/v1/works/:id/kinships", h.ApiPostOne)
	return nil
}
