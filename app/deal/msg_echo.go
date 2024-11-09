package deal

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type handlerEcho struct {
	api API
	ssr msg.Renderer
	log *slog.Logger
}

func newHandlerEcho(a API, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "dealHandlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) ApiPostOne(c echo.Context) error {
	var mto SpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	spec, err := MsgToSpec(mto)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromRoot(root))
}

func (h *handlerEcho) ApiGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := id.ConvertFromString(mto.ID)
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
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := id.ConvertFromString(mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("deal", MsgFromRoot(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Adapter
type kinshipHandlerEcho struct {
	api API
	ssr msg.Renderer
	log *slog.Logger
}

func newKinshipHandlerEcho(a API, r msg.Renderer, l *slog.Logger) *kinshipHandlerEcho {
	name := slog.String("name", "kinshipHandlerEcho")
	return &kinshipHandlerEcho{a, r, l.With(name)}
}

func (h *kinshipHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto KinshipSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	err = mto.Validate()
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

// Adapter
type partHandlerEcho struct {
	api API
	ssr msg.Renderer
	log *slog.Logger
}

func newPartHandlerEcho(a API, r msg.Renderer, l *slog.Logger) *partHandlerEcho {
	name := slog.String("name", "partHandlerEcho")
	return &partHandlerEcho{a, r, l.With(name)}
}

func (h *partHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto PartSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		h.log.Error("mto binding failed", slog.Any("reason", err))
		return err
	}
	err = mto.Validate()
	if err != nil {
		h.log.Error("mto validation failed", slog.Any("reason", err), slog.Any("spec", mto))
		return err
	}
	spec, err := MsgToPartSpec(mto)
	if err != nil {
		h.log.Error("spec mapping failed", slog.Any("reason", err), slog.Any("spec", mto))
		return err
	}
	pe, err := h.api.Involve(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, chnl.MsgFromRoot(pe))
}

// Adapter
type stepHandlerEcho struct {
	api API
	ssr msg.Renderer
	log *slog.Logger
}

func newStepHandlerEcho(a API, r msg.Renderer, l *slog.Logger) *stepHandlerEcho {
	name := slog.String("name", "stepHandlerEcho")
	return &stepHandlerEcho{a, r, l.With(name)}
}

func (h *stepHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto TranSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		h.log.Error("mto binding failed", slog.Any("reason", err))
		return err
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, core.LevelTrace, "transition posting started", slog.Any("mto", mto))
	err = mto.Validate()
	if err != nil {
		h.log.Error("mto validation failed", slog.Any("reason", err), slog.Any("mto", mto))
		return err
	}
	spec, err := MsgToTranSpec(mto)
	if err != nil {
		h.log.Error("spec mapping failed", slog.Any("reason", err), slog.Any("mto", mto))
		return err
	}
	err = h.api.Take(spec)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
