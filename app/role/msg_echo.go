package role

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/id"
)

// Adapter
type handlerEcho struct {
	api API
	log *slog.Logger
}

func newHandlerEcho(a API, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "roleHandlerEcho")
	return &handlerEcho{a, l.With(name)}
}

func (h *handlerEcho) PostOne(c echo.Context) error {
	var dto SpecMsg
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	h.log.Debug("role posting started", slog.Any("dto", dto))
	err = dto.Validate()
	if err != nil {
		return err
	}
	spec, err := MsgToSpec(dto)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromRoot(root))
}

func (h *handlerEcho) GetOne(c echo.Context) error {
	var dto RefMsg
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	ident, err := id.ConvertFromString(dto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.RetrieveLatest(ident)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromRoot(root))
}

func (h *handlerEcho) PatchOne(c echo.Context) error {
	var dto PatchMsg
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	h.log.Debug("role patching started", slog.Any("dto", dto))
	patch, err := MsgToPatch(dto)
	if err != nil {
		return err
	}
	err = h.api.Update(patch)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
