package work

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type workHandlerEcho struct {
	api WorkApi
	ssr msg.Renderer
	log *slog.Logger
}

func newWorkHandlerEcho(a WorkApi, r msg.Renderer, l *slog.Logger) *workHandlerEcho {
	name := slog.String("name", "workHandlerEcho")
	return &workHandlerEcho{a, r, l.With(name)}
}

func (h *workHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto WorkSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	spec, err := MsgToWorkSpec(mto)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromWorkRoot(root))
}

func (h *workHandlerEcho) ApiGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := core.FromString[Work](mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromWorkRoot(root))
}

func (h *workHandlerEcho) SsrGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := core.FromString[Work](mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("work", MsgFromWorkRoot(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Adapter
type kinshipHandlerEcho struct {
	api WorkApi
	ssr msg.Renderer
	log *slog.Logger
}

func newKinshipHandlerEcho(a WorkApi, r msg.Renderer, l *slog.Logger) *kinshipHandlerEcho {
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
