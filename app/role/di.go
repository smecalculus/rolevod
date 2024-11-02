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
		newPresenterEcho,
		fx.Annotate(newRepoPgx, fx.As(new(repo))),
		newKinshipHandlerEcho,
		fx.Annotate(newKinshipRepoPgx, fx.As(new(kinshipRepo))),
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
	),
	fx.Invoke(
		cfgApiEcho,
		cfgSsrEcho,
		cfgKinshipEcho,
	),
)

//go:embed *.html
var formFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("role").Funcs(sprig.FuncMap()).ParseFS(formFs, "*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgApiEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/api/v1/roles", h.PostOne)
	e.GET("/api/v1/roles/:id", h.GetOne)
	e.PUT("/api/v1/roles/:id", h.PutOne)
	return nil
}

func cfgSsrEcho(e *echo.Echo, h *presenterEcho) error {
	e.POST("/ssr/roles", h.PostOne)
	e.GET("/ssr/roles", h.GetMany)
	e.GET("/ssr/roles/:id", h.GetOne)
	return nil
}

func cfgKinshipEcho(e *echo.Echo, h *kinshipHandlerEcho) error {
	e.POST("/api/v1/roles/:id/kinships", h.PostOne)
	return nil
}
