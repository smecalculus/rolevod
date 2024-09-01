package env

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type handlerEcho struct {
	api EnvApi
	ssr msg.Renderer
	log *slog.Logger
}

func newHandlerEcho(a EnvApi, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "env.handlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) ApiPostOne(c echo.Context) error {
	var spec EnvSpec
	err := c.Bind(&spec)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromEnvRoot(root))
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
	return c.JSON(http.StatusOK, MsgFromEnvRoot(root))
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
	html, err := h.ssr.Render("envRoot", MsgFromEnvRoot(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
