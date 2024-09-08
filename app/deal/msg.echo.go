package deal

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type dealHandlerEcho struct {
	api DealApi
	ssr msg.Renderer
	log *slog.Logger
}

func newDealHandlerEcho(a DealApi, r msg.Renderer, l *slog.Logger) *dealHandlerEcho {
	name := slog.String("name", "dealHandlerEcho")
	return &dealHandlerEcho{a, r, l.With(name)}
}

func (h *dealHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto DealSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	spec, err := MsgToDealSpec(mto)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromDealRoot(root))
}

func (h *dealHandlerEcho) ApiGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := id.String[ID](mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromDealRoot(root))
}

func (h *dealHandlerEcho) SsrGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := id.String[ID](mto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("deal", MsgFromDealRoot(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Adapter
type kinshipHandlerEcho struct {
	api DealApi
	ssr msg.Renderer
	log *slog.Logger
}

func newKinshipHandlerEcho(a DealApi, r msg.Renderer, l *slog.Logger) *kinshipHandlerEcho {
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

// Adapter
type partHandlerEcho struct {
	api DealApi
	ssr msg.Renderer
	log *slog.Logger
}

func newPartHandlerEcho(a DealApi, r msg.Renderer, l *slog.Logger) *partHandlerEcho {
	name := slog.String("name", "partHandlerEcho")
	return &partHandlerEcho{a, r, l.With(name)}
}

func (h *partHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto DealSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	spec, err := MsgToDealSpec(mto)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromDealRoot(root))
}

// Adapter
type tranHandlerEcho struct {
	api DealApi
	ssr msg.Renderer
	log *slog.Logger
}

func newTranHandlerEcho(a DealApi, r msg.Renderer, l *slog.Logger) *tranHandlerEcho {
	name := slog.String("name", "tranHandlerEcho")
	return &tranHandlerEcho{a, r, l.With(name)}
}

func (h *tranHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto DealSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	spec, err := MsgToDealSpec(mto)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromDealRoot(root))
}
