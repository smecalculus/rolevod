package role

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type roleHandlerEcho struct {
	api RoleApi
	ssr msg.Renderer
	log *slog.Logger
}

func newRoleHandlerEcho(ra RoleApi, r msg.Renderer, l *slog.Logger) *roleHandlerEcho {
	name := slog.String("name", "roleHandlerEcho")
	return &roleHandlerEcho{ra, r, l.With(name)}
}

func (h *roleHandlerEcho) ApiPostOne(c echo.Context) error {
	var mto RoleSpecMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	h.log.Debug("role posting started", slog.Any("mto", mto))
	spec, err := MsgToRoleSpec(mto)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromRoleRoot(root))
}

func (h *roleHandlerEcho) ApiGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := id.String[ID](mto.ID)
	if err != nil {
		return err
	}
	rr, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromRoleRoot(rr))
}

func (h *roleHandlerEcho) ApiPutOne(c echo.Context) error {
	var mto RoleRootMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	rr, err := MsgToRoleRoot(mto)
	if err != nil {
		return err
	}
	err = h.api.Update(rr)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *roleHandlerEcho) SsrGetOne(c echo.Context) error {
	var mto RefMsg
	err := c.Bind(&mto)
	if err != nil {
		return err
	}
	id, err := id.String[ID](mto.ID)
	if err != nil {
		return err
	}
	rr, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("tp", MsgFromRoleRoot(rr))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Adapter
type kinshipHandlerEcho struct {
	api RoleApi
	ssr msg.Renderer
	log *slog.Logger
}

func newKinshipHandlerEcho(a RoleApi, r msg.Renderer, l *slog.Logger) *kinshipHandlerEcho {
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
