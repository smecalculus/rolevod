package web

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/msg"
)

// adapter
type HandlerEcho struct {
	renderer msg.Renderer
	log      *slog.Logger
}

func NewHandlerEcho(r msg.Renderer, l *slog.Logger) *HandlerEcho {
	name := slog.String("name", "web.HandlerEcho")
	return &HandlerEcho{r, l.With(name)}
}

func (h *HandlerEcho) Home(c echo.Context) error {
	blob, err := h.renderer.Render("home.go.html", nil)
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, blob)
}
