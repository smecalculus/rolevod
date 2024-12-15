package role

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
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
		h.log.Error("dto binding failed")
		return err
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, core.LevelTrace, "role posting started", slog.Any("dto", dto))
	err = dto.Validate()
	if err != nil {
		h.log.Error("dto validation failed")
		return err
	}
	spec, err := MsgToSpec(dto)
	if err != nil {
		h.log.Error("dto mapping failed")
		return err
	}
	snap, err := h.api.Create(spec)
	if err != nil {
		h.log.Error("role creation failed")
		return err
	}
	h.log.Log(ctx, core.LevelTrace, "role posting succeeded", slog.Any("id", snap.ID))
	return c.JSON(http.StatusCreated, MsgFromSnap(snap))
}

func (h *handlerEcho) GetOne(c echo.Context) error {
	var dto IdentMsg
	err := c.Bind(&dto)
	if err != nil {
		h.log.Error("dto binding failed")
		return err
	}
	err = dto.Validate()
	if err != nil {
		h.log.Error("dto validation failed")
		return err
	}
	id, err := id.ConvertFromString(dto.ID)
	if err != nil {
		h.log.Error("dto mapping failed")
		return err
	}
	snap, err := h.api.Retrieve(id)
	if err != nil {
		h.log.Error("root retrieval failed")
		return err
	}
	return c.JSON(http.StatusOK, MsgFromSnap(snap))
}

func (h *handlerEcho) PatchOne(c echo.Context) error {
	var dto SnapMsg
	err := c.Bind(&dto)
	if err != nil {
		h.log.Error("dto binding failed")
		return err
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, core.LevelTrace, "role patching started", slog.Any("dto", dto))
	err = dto.Validate()
	if err != nil {
		h.log.Error("dto validation failed")
		return err
	}
	reqSnap, err := MsgToSnap(dto)
	if err != nil {
		h.log.Error("dto mapping failed")
		return err
	}
	resSnap, err := h.api.Modify(reqSnap)
	if err != nil {
		h.log.Error("role modification failed")
		return err
	}
	h.log.Log(ctx, core.LevelTrace, "role patching succeeded", slog.Any("ref", ConvertSnapToRef(resSnap)))
	return c.JSON(http.StatusOK, MsgFromSnap(resSnap))
}
