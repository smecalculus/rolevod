package env

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// adapter
type repoPgx struct {
	conv dataConverter
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(c dataConverter, p *pgxpool.Pool, log *slog.Logger) *repoPgx {
	name := slog.String("name", "env.repoPgx")
	return &repoPgx{c, p, log.With(name)}
}

func (r *repoPgx) Insert(root Root) error {
	return nil
}
