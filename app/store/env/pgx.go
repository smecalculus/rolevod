package env

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/app/core/env"
)

// adapter
type RepoPgx struct {
	pgx *pgxpool.Pool
	log *slog.Logger
}

func NewRepoPgx(pgx *pgxpool.Pool, log *slog.Logger) *RepoPgx {
	name := slog.String("name", "env.RepoPgx")
	return &RepoPgx{pgx, log.With(name)}
}

func (r *RepoPgx) Insert(er env.Root) error {
	return nil
}
