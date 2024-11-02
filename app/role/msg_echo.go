package role

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

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
	name := slog.String("name", "roleHandlerEcho")
	return &handlerEcho{a, r, l.With(name)}
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
	ident, err := id.StringToID(dto.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(ident)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromRoot(root))
}

func (h *handlerEcho) PutOne(c echo.Context) error {
	var dto RootMsg
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	root, err := MsgToRoot(dto)
	if err != nil {
		return err
	}
	err = h.api.Update(root)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
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

func (h *kinshipHandlerEcho) PostOne(c echo.Context) error {
	var dto KinshipSpecMsg
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	spec, err := MsgToKinshipSpec(dto)
	if err != nil {
		return err
	}
	err = h.api.Establish(spec)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
