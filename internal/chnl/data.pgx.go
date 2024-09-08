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
type repoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	name := slog.String("name", "chnl.repoPgx")
	return &repoPgx{p, l.With(name)}
}

func (r *repoPgx) Insert(root Root) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := DataFromRoot(root)
	query := `
		INSERT INTO channels (
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
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("channel", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *repoPgx) SelectAll() ([]Ref, error) {
	roots := make([]Ref, 5)
	for i := range 5 {
		roots[i] = Ref{ID: id.New[ID](), Name: fmt.Sprintf("Root%v", i)}
	}
	return roots, nil
}

func (r *repoPgx) SelectById(id id.ADT[ID]) (Root, error) {
	return Root{ID: id, Name: "Root"}, nil
}
