//go:build !goverter

package seat

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/seat",
	fx.Provide(
		fx.Annotate(newSeatService, fx.As(new(SeatApi))),
		fx.Annotate(newSeatRepoPgx, fx.As(new(SeatRepo))),
	),
	fx.Provide(
		fx.Private,
		newSeatHandlerEcho,
		newKinshipHandlerEcho,
		fx.Annotate(newKinshipRepoPgx, fx.As(new(kinshipRepo))),
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
	),
	fx.Invoke(
		cfgSeatEcho,
		cfgKinshipEcho,
	),
)

//go:embed all:view
var viesFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("seat").Funcs(sprig.FuncMap()).ParseFS(viesFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgSeatEcho(e *echo.Echo, h *seatHandlerEcho) error {
	e.POST("/api/v1/seats", h.ApiPostOne)
	e.GET("/api/v1/seats/:id", h.ApiGetOne)
	e.GET("/ssr/seats/:id", h.SsrGetOne)
	return nil
}

func cfgKinshipEcho(e *echo.Echo, h *kinshipHandlerEcho) error {
	e.POST("/api/v1/seats/:id/kinships", h.ApiPostOne)
	return nil
}
