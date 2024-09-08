package chnl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/id"
)

// Adapter
type RepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(p *pgxpool.Pool, l *slog.Logger) *RepoPgx {
	name := slog.String("name", "RepoPgx")
	return &RepoPgx{p, l.With(name)}
}

func (r *RepoPgx) Insert(root Root) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := DataFromRoot(root)
	query := `
		INSERT INTO s (
			id, name
		) VALUES (
			@id, @name
		)`
	args := pgx.NamedArgs{
		"id":   dto.ID,
		"name": dto.Name,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *RepoPgx) SelectAll() ([]Ref, error) {
	roots := make([]Ref, 5)
	for i := range 5 {
		roots[i] = Ref{ID: id.New[ID](), Name: fmt.Sprintf("Root%v", i)}
	}
	return roots, nil
}

func (r *RepoPgx) SelectById(id id.ADT[ID]) (Root, error) {
	return Root{ID: id, Name: "Root"}, nil
}
