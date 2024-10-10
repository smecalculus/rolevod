package step

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/lib/core"
)

// Adapter
type repoPgx[T Root] struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx[T Root](p *pgxpool.Pool, l *slog.Logger) *repoPgx[T] {
	name := slog.String("name", "stepRepoPgx[T]")
	return &repoPgx[T]{p, l.With(name)}
}

// for compilation purposes
func newRepo[T Root]() Repo[T] {
	return &repoPgx[T]{}
}

func (r *repoPgx[T]) Insert(root Root) error {
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
			id, kind, pid, vid, ctx, term
		) VALUES (
			@id, @kind, @pid, @vid, @ctx, @term
		)`
	args := pgx.NamedArgs{
		"id":   dto.ID,
		"kind": dto.K,
		"pid":  dto.PID,
		"vid":  dto.VID,
		"ctx":  dto.Ctx,
		"term": dto.Term,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *repoPgx[T]) SelectAll() ([]Ref, error) {
	return nil, nil
}

func (r *repoPgx[T]) SelectByID(rid ID) (*T, error) {
	query := `
		SELECT
			id, kind, pid, vid, ctx term
		FROM steps
		WHERE id = $1`
	return r.execute(query, rid.String())
}

func (r *repoPgx[T]) SelectByPID(pid chnl.ID) (*T, error) {
	query := `
		SELECT
			id, kind, pid, vid, ctx, term
		FROM steps
		WHERE pid = $1`
	return r.execute(query, pid.String())
}

func (r *repoPgx[T]) SelectByVID(vid chnl.ID) (*T, error) {
	query := `
		SELECT
			id, kind, pid, vid, ctx, term
		FROM steps
		WHERE vid = $1`
	return r.execute(query, vid.String())
}

func (r *repoPgx[T]) execute(query string, arg string) (*T, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, arg)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return nil, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return nil, err
	}
	generic, err := dataToRoot(&dto)
	if err != nil {
		r.log.Error("dto mapping failed", slog.Any("reason", err))
		return nil, err
	}
	concrete, ok := generic.(T)
	if !ok {
		err = ErrUnexpectedRootType(generic)
		r.log.Error("step selection failed", slog.Any("reason", err), slog.Any("dto", dto))
		return nil, err
	}
	r.log.Log(ctx, core.LevelTrace, "step selection succeeded", slog.Any("root", concrete))
	return &concrete, nil
}
