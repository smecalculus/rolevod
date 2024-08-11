package env

import (
	"fmt"
	"log/slog"
	"mime"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"

	env "smecalculus/rolevod/app/env/core"
)

// adapter
type HandlerEcho struct {
	api  env.Api
	html msg.Renderer
	log  *slog.Logger
}

func NewHandlerEcho(a env.Api, r msg.Renderer, l *slog.Logger) *HandlerEcho {
	name := slog.String("name", "env.HandlerEcho")
	return &HandlerEcho{a, r, l.With(name)}
}

func (h *HandlerEcho) Post(c echo.Context) error {
	var spec env.Spec
	err1 := c.Bind(&spec)
	if err1 != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	_, err2 := h.api.Create(spec)
	if err2 != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	return c.NoContent(http.StatusOK)
}

func (h *HandlerEcho) Get(c echo.Context) error {
	mediaType, _, err := mime.ParseMediaType(c.Request().Header.Get(echo.HeaderAccept))
	if err != nil {
		return err
	}
	switch mediaType {
	case echo.MIMEApplicationJSON, echo.MIMETextPlain, msg.MIMEAnyAny:
		return c.JSON(http.StatusOK, env.Env{Id: core.New[env.Id]()})
	case echo.MIMETextHTML, echo.MIMETextHTMLCharsetUTF8:
		blob, err := h.html.Render("env", env.Env{Id: core.New[env.Id]()})
		if err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, blob)
	default:
		return echo.NewHTTPError(http.StatusUnsupportedMediaType,
			fmt.Sprintf("unsupported media type: %v", mediaType))
	}
}
