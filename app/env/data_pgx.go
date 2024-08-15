package env

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/decl"
)

// adapter
type repoPgx struct {
	conv dataConverter
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(c dataConverter, p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	name := slog.String("name", "env.repoPgx")
	return &repoPgx{c, p, l.With(name)}
}

func (r *repoPgx) Insert(root Root) error {
	return nil
}

func (r *repoPgx) SelectById(id core.ID[Env]) (Root, error) {
	decls := make([]decl.Root, 5)
	for i := range 5 {
		decls[i] = decl.Root{
			ID:   core.New[decl.Decl](),
			Name: fmt.Sprintf("Foo%v", i)}
	}
	return Root{id, "Foo", decls}, nil
}

func (r *repoPgx) SelectAll() ([]Root, error) {
	roots := make([]Root, 5)
	for i := range 5 {
		roots[i] = Root{core.New[Env](), fmt.Sprintf("Foo%v", i), []decl.Root{}}
	}
	return roots, nil
}
