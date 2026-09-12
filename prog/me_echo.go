package prog

import (
	"io"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	sdkprog "github.com/orglang/go-sdk/prog"
)

type controllerEcho struct {
	api API
	log *slog.Logger
}

func newControllerEcho(a API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{a, l.With(name)}
}

func cfgControllerEcho(e *echo.Echo, h *controllerEcho) error {
	e.POST("/api/v1/progs", h.PostSpec)
	return nil
}

func (h *controllerEcho) PostSpec(ctx echo.Context) error {
	text, readErr := io.ReadAll(ctx.Request().Body)
	if readErr != nil {
		h.log.Error("reading failed", slog.Any("err", readErr))
		return readErr
	}
	dto, parseErr := sdkprog.MsgFromText(string(text))
	if parseErr != nil {
		h.log.Error("parsing failed", slog.Any("dto", reflect.TypeFor[sdkprog.Spec]()))
		return parseErr
	}
	spec, convErr := MsgToSpec(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	apiErr := h.api.Create(spec)
	if apiErr != nil {
		return apiErr
	}
	return ctx.NoContent(http.StatusCreated)
}
