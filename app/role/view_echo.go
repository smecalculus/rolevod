package role

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/msg"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/state"
)

// Adapter
type presenterEcho struct {
	api API
	ssr msg.Renderer
	log *slog.Logger
}

func newPresenterEcho(a API, r msg.Renderer, l *slog.Logger) *presenterEcho {
	name := slog.String("name", "rolePresenterEcho")
	return &presenterEcho{a, r, l.With(name)}
}

func (p *presenterEcho) PostOne(c echo.Context) error {
	var dto SpecView
	err := c.Bind(&dto)
	if err != nil {
		p.log.Error("dto binding failed")
		return err
	}
	p.log.Log(c.Request().Context(), core.LevelTrace, "role posting started", slog.Any("dto", dto))
	err = dto.Validate()
	if err != nil {
		p.log.Error("dto validation failed")
		return err
	}
	fqn := sym.CovertFromString(dto.NS).New(dto.Name)
	root, err := p.api.Create(Spec{FQN: fqn, State: state.OneSpec{}})
	if err != nil {
		p.log.Error("root creation failed")
		return err
	}
	html, err := p.ssr.Render("view-one", ViewFromRoot(root))
	if err != nil {
		p.log.Error("view rendering failed")
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *presenterEcho) GetMany(c echo.Context) error {
	refs, err := p.api.RetreiveAll()
	if err != nil {
		p.log.Error("refs retrieving failed")
		return err
	}
	html, err := p.ssr.Render("view-many", ViewFromRefs(refs))
	if err != nil {
		p.log.Error("view rendering failed")
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *presenterEcho) GetOne(c echo.Context) error {
	var dto RefMsg
	err := c.Bind(&dto)
	if err != nil {
		p.log.Error("dto binding failed")
		return err
	}
	ident, err := id.ConvertFromString(dto.ID)
	if err != nil {
		p.log.Error("ident mapping failed")
		return err
	}
	root, err := p.api.Retrieve(ident)
	if err != nil {
		p.log.Error("root retrieving failed")
		return err
	}
	html, err := p.ssr.Render("view-one", ViewFromRoot(root))
	if err != nil {
		p.log.Error("view rendering failed")
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
