package web

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/msg"

	"smecalculus/rolevod/app/role"
)

// Adapter
type handlerEcho struct {
	api role.API
	ssr msg.Renderer
	log *slog.Logger
}

func newHandlerEcho(a role.API, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "webHandlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) Home(c echo.Context) error {
	refs, err := h.api.RetreiveRefs()
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("home.html", role.MsgFromRefs(refs))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
