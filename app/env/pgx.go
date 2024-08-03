package env

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// adapter
type repoPgx struct {
	pgx *pgxpool.Pool
}

func (r *repoPgx) Insert(er envRoot) error {
	return nil
}
