package web

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/msg"

	"smecalculus/rolevod/app/env"
)

// Adapter
type handlerEcho struct {
	api  env.EnvApi
	view msg.Renderer
	log  *slog.Logger
}

func hewHandlerEcho(a env.EnvApi, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "web.handlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) Home(c echo.Context) error {
	roots, err := h.api.RetreiveAll()
	if err != nil {
		return err
	}
	html, err := h.view.Render("home.html", env.MsgFromEnvRoots(roots))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
