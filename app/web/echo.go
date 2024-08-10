package web

import (
	"log/slog"
	"net/http"

	"smecalculus/rolevod/lib/msg"

	"github.com/labstack/echo/v4"
)

// adapter
type handlerEcho struct {
	renderer msg.Renderer
	log      *slog.Logger
}

func (h *handlerEcho) home(c echo.Context) error {
	blob, err := h.renderer.Render("home.go.tmpl", nil)
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, blob)
}
