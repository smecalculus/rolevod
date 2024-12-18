//go:build !goverter

package deal

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/deal",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		newHandlerEcho,
		fx.Annotate(newRepoPgx, fx.As(new(repo))),
		newKinshipHandlerEcho,
		fx.Annotate(newKinshipRepoPgx, fx.As(new(kinshipRepo))),
		newPartHandlerEcho,
		newStepHandlerEcho,
	),
	fx.Invoke(
		cfgDealEcho,
		cfgKinshipEcho,
		cfgPartEcho,
		cfgStepEcho,
	),
)

//go:embed all:view
var viesFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("deal").Funcs(sprig.FuncMap()).ParseFS(viesFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgDealEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/api/v1/deals", h.ApiPostOne)
	e.GET("/api/v1/deals/:id", h.ApiGetOne)
	e.GET("/ssr/deals/:id", h.SsrGetOne)
	return nil
}

func cfgKinshipEcho(e *echo.Echo, h *kinshipHandlerEcho) error {
	e.POST("/api/v1/deals/:id/kinships", h.ApiPostOne)
	return nil
}

func cfgPartEcho(e *echo.Echo, h *partHandlerEcho) error {
	e.POST("/api/v1/deals/:id/parts", h.ApiPostOne)
	return nil
}

func cfgStepEcho(e *echo.Echo, h *stepHandlerEcho) error {
	e.POST("/api/v1/deals/:id/steps", h.ApiPostOne)
	return nil
}
