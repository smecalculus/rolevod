package env

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type repoPgx struct {
	pgx *pgxpool.Pool
}

func (r *repoPgx) Insert(er envRoot) error {
	return nil
}
