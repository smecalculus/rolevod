//go:build !goverter

package env

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/env",
	fx.Provide(
		fx.Annotate(NewService, fx.As(new(Api))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		fx.Annotate(newMsgConverter, fx.As(new(msgConverter))),
		newHandlerEcho,
		fx.Annotate(newDataConverter, fx.As(new(dataConverter))),
		fx.Annotate(newRepoPgx, fx.As(new(Repo))),
	),
	fx.Invoke(
		cfgEcho,
	),
)

//go:embed *.html
var envFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.ParseFS(envFs, "*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func newMsgConverter() msgConverter {
	return &msgConverterImpl{}
}

func newDataConverter() dataConverter {
	return &dataConverterImpl{}
}

func cfgEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/envs", h.Post)
	e.GET("/envs/:id", h.Get)
	return nil
}
