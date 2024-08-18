package env

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/dcl"
)

// Adapter
type repoPgx struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	name := slog.String("name", "env.repoPgx")
	return &repoPgx{p, l.With(name)}
}

func (r *repoPgx) Insert(root AR) error {
	return nil
}

func (r *repoPgx) SelectById(id core.ID[AR]) (AR, error) {
	tpDefs := make([]dcl.TpRoot, 5)
	for i := range 5 {
		tpDefs[i] = dcl.TpRoot{
			ID:   core.New[dcl.AR](),
			Name: fmt.Sprintf("TpRoot%v", i)}
	}
	expDecs := make([]dcl.ExpRoot, 5)
	for i := range 5 {
		expDecs[i] = dcl.ExpRoot{
			ID:   core.New[dcl.AR](),
			Name: fmt.Sprintf("ExpRoot%v", i)}
	}
	return AR{id, "Foo", tpDefs, expDecs}, nil
}

func (r *repoPgx) SelectAll() ([]AR, error) {
	roots := make([]AR, 5)
	for i := range 5 {
		roots[i] = AR{core.New[AR](), fmt.Sprintf("Foo%v", i), []dcl.TpRoot{}, []dcl.ExpRoot{}}
	}
	return roots, nil
}
