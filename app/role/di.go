//go:build !goverter

package role

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/role",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
	),
	fx.Provide(
		fx.Private,
		newHandlerEcho,
		fx.Annotate(newRepoPgx, fx.As(new(repo))),
		newKinshipHandlerEcho,
		fx.Annotate(newKinshipRepoPgx, fx.As(new(kinshipRepo))),
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
	),
	fx.Invoke(
		cfgRoleEcho,
		cfgKinshipEcho,
	),
)

//go:embed all:view
var viewFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("role").Funcs(sprig.FuncMap()).ParseFS(viewFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgRoleEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/api/v1/roles", h.ApiPostOne)
	e.GET("/api/v1/roles/:id", h.ApiGetOne)
	e.PUT("/api/v1/roles/:id", h.ApiPutOne)
	e.GET("/ssr/roles/:id", h.SsrGetOne)
	return nil
}

func cfgKinshipEcho(e *echo.Echo, h *kinshipHandlerEcho) error {
	e.POST("/api/v1/roles/:id/kinships", h.ApiPostOne)
	return nil
}
