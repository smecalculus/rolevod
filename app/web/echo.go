package web

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/msg"
)

// adapter
type handlerEcho struct {
	renderer msg.Renderer
	log      *slog.Logger
}

func newHandlerEcho(r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "web.handlerEcho")
	return &handlerEcho{r, l.With(name)}
}

func (h *handlerEcho) home(c echo.Context) error {
	blob, err := h.renderer.Render("home.go.html", nil)
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, blob)
}
