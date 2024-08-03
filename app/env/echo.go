package env

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// adapter
type handlerEcho struct {
	api Api
}

func (h *handlerEcho) post(c echo.Context) error {
	var es EnvSpec
	err1 := c.Bind(&es)
	if err1 != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	_, err2 := h.api.Create(es)
	if err2 != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	return c.NoContent(http.StatusOK)
}

func (h *handlerEcho) get(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
