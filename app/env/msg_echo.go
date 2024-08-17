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
	conv MsgConverter
	view msg.Renderer
	log  *slog.Logger
}

func newHandlerEcho(a Api, c MsgConverter, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "env.handlerEcho")
	return &handlerEcho{a, c, r, l.With(name)}
}

func (h *handlerEcho) ApiPostOne(c echo.Context) error {
	var spec Spec
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
		return c.JSON(http.StatusOK, h.conv.ToRootMsg(root))
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
	var params GetMsg
	err := c.Bind(&params)
	if err != nil {
		return err
	}
	id, err := core.FromString[Env](params.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, h.conv.ToRootMsg(root))
}

func (h *handlerEcho) SsrGetOne(c echo.Context) error {
	var params GetMsg
	err := c.Bind(&params)
	if err != nil {
		return err
	}
	id, err := core.FromString[Env](params.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.view.Render("envRoot", h.conv.ToRootMsg(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
