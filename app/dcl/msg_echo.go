package dcl

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
)

// adapter
type handlerEcho struct {
	api  Api
	view msg.Renderer
	log  *slog.Logger
}

func newHandlerEcho(a Api, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "decl.handlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) SsrGetOne(c echo.Context) error {
	var params GetMsg
	err := c.Bind(&params)
	if err != nil {
		return err
	}
	id, err := core.FromString[Dcl](params.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	var html []byte
	switch decl := root.(type) {
	case TpDef:
		html, err = h.view.Render("declRoot", ToRootMsg(decl))
	}
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
