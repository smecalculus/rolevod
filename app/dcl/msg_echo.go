package dcl

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type tpHandlerEcho struct {
	api  TpApi
	view msg.Renderer
	log  *slog.Logger
}

func newTpHandlerEcho(a TpApi, r msg.Renderer, l *slog.Logger) *tpHandlerEcho {
	name := slog.String("name", "dcl.tpHandlerEcho")
	return &tpHandlerEcho{a, r, l.With(name)}
}

func (h *tpHandlerEcho) SsrGetOne(c echo.Context) error {
	var ref RefMsg
	err := c.Bind(&ref)
	if err != nil {
		return err
	}
	id, err := core.FromString[AR](ref.ID)
	if err != nil {
		return err
	}
	ar, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.view.Render("dclRoot", MsgFromTpRoot(ar))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Adapter
type expHandlerEcho struct {
	api  ExpApi
	view msg.Renderer
	log  *slog.Logger
}

func newExpHandlerEcho(a ExpApi, r msg.Renderer, l *slog.Logger) *expHandlerEcho {
	name := slog.String("name", "dcl.expHandlerEcho")
	return &expHandlerEcho{a, r, l.With(name)}
}

func (h *expHandlerEcho) SsrGetOne(c echo.Context) error {
	var ref RefMsg
	err := c.Bind(&ref)
	if err != nil {
		return err
	}
	id, err := core.FromString[AR](ref.ID)
	if err != nil {
		return err
	}
	ar, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.view.Render("dclRoot", MsgFromExpRoot(ar))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}