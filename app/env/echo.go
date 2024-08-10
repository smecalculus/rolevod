package env

import (
	"fmt"
	"log/slog"
	"mime"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/msg"
)

// adapter
type handlerEcho struct {
	env  Api
	html msg.Renderer
	log  *slog.Logger
}

func (h *handlerEcho) post(c echo.Context) error {
	var es EnvSpec
	err1 := c.Bind(&es)
	if err1 != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	_, err2 := h.env.Create(es)
	if err2 != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	return c.NoContent(http.StatusOK)
}

func (h *handlerEcho) get(c echo.Context) error {
	acceptValue := c.Request().Header.Get(echo.HeaderAccept)
	if acceptValue == "" {
		return c.JSON(http.StatusOK, Env{"foo"})
	}
	mediaType, _, err := mime.ParseMediaType(acceptValue)
	if err != nil {
		return err
	}
	switch mediaType {
	case "*/*", echo.MIMEApplicationJSON:
		return c.JSON(http.StatusOK, Env{"foo"})
	case echo.MIMETextHTML, echo.MIMETextHTMLCharsetUTF8:
		blob, err := h.html.Render("env", Env{"foo"})
		if err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, blob)
	default:
		return c.String(http.StatusUnsupportedMediaType,
			fmt.Sprintf("unsupported accept media type: %v", mediaType))
	}
}
