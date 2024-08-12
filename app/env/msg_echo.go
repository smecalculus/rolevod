package env

import (
	"fmt"
	"log/slog"
	"mime"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
)

// adapter
type handlerEcho struct {
	api  Api
	json msgConverter
	html msg.Renderer
	log  *slog.Logger
}

func newHandlerEcho(a Api, c msgConverter, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "env.handlerEcho")
	return &handlerEcho{a, c, r, l.With(name)}
}

func (h *handlerEcho) Post(c echo.Context) error {
	var spec Spec
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

func (h *handlerEcho) Get(c echo.Context) error {
	mediaType, _, err := mime.ParseMediaType(c.Request().Header.Get(echo.HeaderAccept))
	if err != nil {
		return err
	}
	switch mediaType {
	case echo.MIMEApplicationJSON, echo.MIMETextPlain, msg.MIMEAnyAny:
		return c.JSON(http.StatusOK, Root{Id: core.New[Id]()})
	case echo.MIMETextHTML, echo.MIMETextHTMLCharsetUTF8:
		blob, err := h.html.Render("env", Root{Id: core.New[Id]()})
		if err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, blob)
	default:
		return echo.NewHTTPError(http.StatusUnsupportedMediaType,
			fmt.Sprintf("unsupported media type: %v", mediaType))
	}
}
