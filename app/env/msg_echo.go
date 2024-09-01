package env

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/msg"
)

// Adapter
type envHandlerEcho struct {
	api EnvApi
	ssr msg.Renderer
	log *slog.Logger
}

func newEnvHandlerEcho(a EnvApi, r msg.Renderer, l *slog.Logger) *envHandlerEcho {
	name := slog.String("name", "ws.envHandlerEcho")
	return &envHandlerEcho{a, r, l.With(name)}
}

func (h *envHandlerEcho) ApiPostOne(c echo.Context) error {
	var spec EnvSpec
	err := c.Bind(&spec)
	if err != nil {
		return err
	}
	root, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromEnvRoot(root))
}

func (h *envHandlerEcho) ApiGetOne(c echo.Context) error {
	var ref RefMsg
	err := c.Bind(&ref)
	if err != nil {
		return err
	}
	id, err := core.FromString[AR](ref.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromEnvRoot(root))
}

func (h *envHandlerEcho) SsrGetOne(c echo.Context) error {
	var ref RefMsg
	err := c.Bind(&ref)
	if err != nil {
		return err
	}
	id, err := core.FromString[AR](ref.ID)
	if err != nil {
		return err
	}
	root, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("envRoot", MsgFromEnvRoot(root))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Adapter
type introHandlerEcho struct {
	api EnvApi
	log *slog.Logger
}

func newIntroHandlerEcho(a EnvApi, l *slog.Logger) *introHandlerEcho {
	name := slog.String("name", "ws.introHandlerEcho")
	return &introHandlerEcho{a, l.With(name)}
}

func (h *introHandlerEcho) ApiPostOne(c echo.Context) error {
	var msg IntroMsg
	err := c.Bind(&msg)
	if err != nil {
		return err
	}
	intro, err := MsgToIntro(msg)
	if err != nil {
		return err
	}
	err = h.api.Introduce(intro)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
