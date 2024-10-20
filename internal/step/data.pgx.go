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
type repoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	name := slog.String("name", "stepRepoPgx")
	return &repoPgx{p, l.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return &repoPgx{}
}

func (r *repoPgx) Insert(root Root) error {
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
			id, kind, pid, vid, term
		) VALUES (
			@id, @kind, @pid, @vid, @term
		)`
	args := pgx.NamedArgs{
		"id":   dto.ID,
		"kind": dto.K,
		"pid":  dto.PID,
		"vid":  dto.VID,
		"term": dto.Term,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *repoPgx) SelectAll() ([]Ref, error) {
	return nil, nil
}

func (r *repoPgx) SelectByID(rid ID) (Root, error) {
	query := `
		SELECT
			id, kind, pid, vid, ctx term
		FROM steps
		WHERE id = $1`
	return r.execute(query, rid.String())
}

func (r *repoPgx) SelectByPID(pid chnl.ID) (Root, error) {
	query := `
		SELECT
			id, kind, pid, vid, term
		FROM steps
		WHERE pid = $1`
	return r.execute(query, pid.String())
}

func (r *repoPgx) SelectByVID(vid chnl.ID) (Root, error) {
	query := `
		SELECT
			id, kind, pid, vid, term
		FROM steps
		WHERE vid = $1`
	return r.execute(query, vid.String())
}

func (r *repoPgx) execute(query string, arg string) (Root, error) {
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
	root, err := dataToRoot(&dto)
	if err != nil {
		r.log.Error("dto mapping failed", slog.Any("reason", err))
		return nil, err
	}
	r.log.Log(ctx, core.LevelTrace, "step selection succeeded", slog.Any("root", root))
	return root, nil
}
