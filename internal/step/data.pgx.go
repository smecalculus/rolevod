package step

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

// Adapter
type repoPgx[T root] struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx[T root](p *pgxpool.Pool, l *slog.Logger) *repoPgx[T] {
	name := slog.String("name", "step.repoPgx[T]")
	return &repoPgx[T]{p, l.With(name)}
}

// for compilation purposes
func newRepo[T root]() Repo[T] {
	return &repoPgx[T]{}
}

func (r *repoPgx[T]) Insert(root root) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto, err := dataFromRoot(root)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO steps (
			id, pre_id, via_id, payload
		) VALUES (
			@id, @pre_id, @via_id, @payload
		)`
	args := pgx.NamedArgs{
		"id":      dto.ID,
		"pre_id":  dto.PreID,
		"via_id":  dto.ViaID,
		"payload": dto.Payload,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("step", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *repoPgx[T]) SelectAll() ([]Ref, error) {
	query := `
		SELECT
			id
		FROM steps`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[refData])
	if err != nil {
		return nil, err
	}
	return DataToRefs(dtos)
}

func (r *repoPgx[T]) SelectByID(rid id.ADT[ID]) (*T, error) {
	query := `
		SELECT
			id, pre_id, via_id, payload
		FROM steps
		WHERE id=$1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
	if err != nil {
		return nil, err
	}
	generic, err := dataToRoot(&dto)
	if err != nil {
		return nil, err
	}
	concrete, ok := generic.(T)
	if !ok {
		return nil, ErrUnexpectedStep
	}
	return &concrete, nil
}

func (r *repoPgx[T]) SelectByChID(vid id.ADT[chnl.ID]) (*T, error) {
	query := `
		SELECT
			id, pre_id, via_id, payload
		FROM steps
		WHERE via_id=$1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, vid.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
	if err != nil {
		return nil, err
	}
	generic, err := dataToRoot(&dto)
	if err != nil {
		return nil, err
	}
	concrete, ok := generic.(T)
	if !ok {
		return nil, ErrUnexpectedStep
	}
	return &concrete, nil
}
