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

// Adapter
type handlerEcho struct {
	api  Api
	view msg.Renderer
	log  *slog.Logger
}

func newHandlerEcho(a Api, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "env.handlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) ApiPostOne(c echo.Context) error {
	var spec AS
	err := c.Bind(&spec)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	mediaType, _, err := mime.ParseMediaType(c.Request().Header.Get(echo.HeaderAccept))
	if err != nil {
		return err
	}
	switch mediaType {
	case echo.MIMEApplicationJSON, echo.MIMETextPlain, msg.MIMEAnyAny:
		return c.JSON(http.StatusOK, MsgFromRoot(root))
	case echo.MIMETextHTML, echo.MIMETextHTMLCharsetUTF8:
		html, err := h.view.Render("root", root)
		if err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, html)
	default:
		return echo.NewHTTPError(http.StatusUnsupportedMediaType,
			fmt.Sprintf("unsupported media type: %v", mediaType))
	}
}

func (h *handlerEcho) ApiGetOne(c echo.Context) error {
	var ref RefMsg
	err := c.Bind(&ref)
	if err != nil {
		return err
	}
	id, err := core.FromString[AR](ref.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromRoot(root))
}

func (h *handlerEcho) SsrGetOne(c echo.Context) error {
	var ref RefMsg
	err := c.Bind(&ref)
	if err != nil {
		return err
	}
	id, err := core.FromString[AR](ref.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.view.Render("envRoot", MsgFromRoot(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
