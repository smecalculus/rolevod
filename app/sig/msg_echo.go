package sig

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type sigHandlerEcho struct {
	api Api
	ssr msg.Renderer
	log *slog.Logger
}

func newSigHandlerEcho(a Api, r msg.Renderer, l *slog.Logger) *sigHandlerEcho {
	name := slog.String("name", "sigHandlerEcho")
	return &sigHandlerEcho{a, r, l.With(name)}
}

func (h *sigHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto SigSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		h.log.Error("mto binding failed", slog.Any("reason", err))
		return err
	}
	err = mto.Validate()
	if err != nil {
		h.log.Error("mto validation failed", slog.Any("reason", err), slog.Any("mto", mto))
		return err
	}
	spec, err := MsgToSigSpec(mto)
	if err != nil {
		h.log.Error("mto conversion failed", slog.Any("reason", err), slog.Any("mto", mto))
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromSigRoot(root))
}

func (h *sigHandlerEcho) ApiGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := id.StringToID(mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromSigRoot(root))
}

func (h *sigHandlerEcho) SsrGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := id.StringToID(mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("sig", MsgFromSigRoot(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Adapter
type kinshipHandlerEcho struct {
	api Api
	ssr msg.Renderer
	log *slog.Logger
}

func newKinshipHandlerEcho(a Api, r msg.Renderer, l *slog.Logger) *kinshipHandlerEcho {
	name := slog.String("name", "kinshipHandlerEcho")
	return &kinshipHandlerEcho{a, r, l.With(name)}
}

func (h *kinshipHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto KinshipSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	spec, err := MsgToKinshipSpec(mto)
	if err != nil {
		return err
	}
	err = h.api.Establish(spec)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
