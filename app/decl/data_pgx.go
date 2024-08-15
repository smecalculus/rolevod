package decl

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"
)

// adapter
type repoPgx struct {
	conv dataConverter
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(c dataConverter, p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	name := slog.String("name", "decl.repoPgx")
	return &repoPgx{c, p, l.With(name)}
}

func (r *repoPgx) Insert(root Root) error {
	return nil
}

func (r *repoPgx) SelectById(id core.ID[Decl]) (Root, error) {
	return Root{id, "Bar"}, nil
}

func (r *repoPgx) SelectAll() ([]Root, error) {
	roots := make([]Root, 5)
	for i := range 5 {
		roots[i] = Root{core.New[Decl](), fmt.Sprintf("Bar%v", i)}
	}
	return roots, nil
}
