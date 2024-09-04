package web

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/msg"

	"smecalculus/rolevod/app/ws"
)

// Adapter
type handlerEcho struct {
	api  ws.EnvApi
	ssr msg.Renderer
	log  *slog.Logger
}

func newHandlerEcho(a ws.EnvApi, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "web.handlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) Home(c echo.Context) error {
	roots, err := h.api.RetreiveAll()
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("home.html", ws.MsgFromEnvRoots(roots))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
