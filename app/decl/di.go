//go:build !goverter

package decl

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/decl",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(Api))),
		fx.Annotate(newMsgConverter, fx.As(new(MsgConverter))),
	),
	fx.Provide(
		fx.Private,
		newHandlerEcho,
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
		fx.Annotate(newDataConverter, fx.As(new(dataConverter))),
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

func newMsgConverter() MsgConverter {
	// return nil
	return &MsgConverterImpl{}
}

func newDataConverter() dataConverter {
	// return nil
	return &dataConverterImpl{}
}

func cfgEcho(e *echo.Echo, h *handlerEcho) error {
	e.GET("/gui/decls/:id", h.GuiGetOne)
	return nil
}
