package dcl

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"
)

// adapter
type repoPgx struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	name := slog.String("name", "decl.repoPgx")
	return &repoPgx{p, l.With(name)}
}

func (r *repoPgx) Insert(tpDef TpDef) error {
	return nil
}

func (r *repoPgx) SelectById(id core.ID[Dcl]) (TpDef, error) {
	return TpDef{id, "TpDef"}, nil
}

func (r *repoPgx) SelectAll() ([]TpDef, error) {
	tpDefs := make([]TpDef, 5)
	for i := range 5 {
		tpDefs[i] = TpDef{core.New[Dcl](), fmt.Sprintf("TpDef%v", i)}
	}
	return tpDefs, nil
}
