package env

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// adapter
type repoPgx struct {
	pgx    *pgxpool.Pool
	logger *slog.Logger
}

func (r *repoPgx) Insert(er envRoot) error {
	return nil
}
