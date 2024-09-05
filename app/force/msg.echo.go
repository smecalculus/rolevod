package force

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type forceHandlerEcho struct {
	api ForceApi
	ssr msg.Renderer
	log *slog.Logger
}

func newForceHandlerEcho(a ForceApi, r msg.Renderer, l *slog.Logger) *forceHandlerEcho {
	name := slog.String("name", "forceHandlerEcho")
	return &forceHandlerEcho{a, r, l.With(name)}
}

func (h *forceHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto ForceSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	spec, err := MsgToForceSpec(mto)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromForceRoot(root))
}

func (h *forceHandlerEcho) ApiGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := core.FromString[Force](mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromForceRoot(root))
}

func (h *forceHandlerEcho) SsrGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := core.FromString[Force](mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("force", MsgFromForceRoot(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Adapter
type kinshipHandlerEcho struct {
	api ForceApi
	ssr msg.Renderer
	log *slog.Logger
}

func newKinshipHandlerEcho(a ForceApi, r msg.Renderer, l *slog.Logger) *kinshipHandlerEcho {
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
