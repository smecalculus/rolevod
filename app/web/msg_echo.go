package web

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/msg"

	"smecalculus/rolevod/app/env"
)

// adapter
type handlerEcho struct {
	api  env.Api
	conv env.MsgConverter
	view msg.Renderer
	log  *slog.Logger
}

func hewHandlerEcho(a env.Api, c env.MsgConverter, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "web.handlerEcho")
	return &handlerEcho{a, c, r, l.With(name)}
}

func (h *handlerEcho) Home(c echo.Context) error {
	roots, err := h.api.RetreiveAll()
	if err != nil {
		return err
	}
	blob, err := h.view.Render("home.html", h.conv.ToRootMsgs(roots))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, blob)
}
