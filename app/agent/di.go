//go:build !goverter

package agent

import (
	"embed"
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/msg"
)

var Module = fx.Module("app/agent",
	fx.Provide(
		fx.Annotate(newAgentService, fx.As(new(AgentApi))),
	),
	fx.Provide(
		fx.Private,
		newAgentHandlerEcho,
		fx.Annotate(newAgentRepoPgx, fx.As(new(agentRepo))),
		newKinshipHandlerEcho,
		fx.Annotate(newKinshipRepoPgx, fx.As(new(kinshipRepo))),
		fx.Annotate(newRenderer, fx.As(new(msg.Renderer))),
	),
	fx.Invoke(
		cfgAgentEcho,
		cfgKinshipEcho,
	),
)

//go:embed all:view
var viesFs embed.FS

func newRenderer(l *slog.Logger) (*msg.RendererStdlib, error) {
	t, err := template.New("agent").Funcs(sprig.FuncMap()).ParseFS(viesFs, "*/*.html")
	if err != nil {
		return nil, err
	}
	return msg.NewRendererStdlib(t, l), nil
}

func cfgAgentEcho(e *echo.Echo, h *agentHandlerEcho) error {
	e.POST("/api/v1/agents", h.ApiPostOne)
	e.GET("/api/v1/agents/:id", h.ApiGetOne)
	e.GET("/ssr/agents/:id", h.SsrGetOne)
	return nil
}

func cfgKinshipEcho(e *echo.Echo, h *kinshipHandlerEcho) error {
	e.POST("/api/v1/agents/:id/kinships", h.ApiPostOne)
	return nil
}
