package team

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

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
	name := slog.String("name", "teamHandlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) PostOne(c echo.Context) error {
	var dto SpecMsg
	err := c.Bind(&dto)
	if err != nil {
		h.log.Error("dto binding failed", slog.Any("reason", err))
		return err
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, core.LevelTrace, "posting started", slog.Any("dto", dto))
	err = dto.Validate()
	if err != nil {
		h.log.Error("dto validation failed", slog.Any("reason", err), slog.Any("dto", dto))
		return err
	}
	spec, err := MsgToSpec(dto)
	if err != nil {
		h.log.Error("dto conversion failed", slog.Any("reason", err), slog.Any("dto", dto))
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromRoot(root))
}

func (h *handlerEcho) GetOne(c echo.Context) error {
	var dto IdentMsg
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	id, err := id.ConvertFromString(dto.ID)
	if err != nil {
		return err
	}
	snap, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromSnap(snap))
}
