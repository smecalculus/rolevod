package pool

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
)

// Adapter
type repoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newRepoPgx(p *pgxpool.Pool, l *slog.Logger) *repoPgx {
	name := slog.String("name", "poolRepoPgx")
	return &repoPgx{p, l.With(name)}
}

func (r *repoPgx) Insert(root Root) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := DataFromRoot(root)
	insertRoot := `
		insert into pool_roots (
			pool_id, rev, title, sup_id
		) values (
			@pool_id, @rev, @title, @sup_id
		)`
	rootArgs := pgx.NamedArgs{
		"pool_id": dto.ID,
		"rev":     dto.Rev,
		"title":   dto.Title,
		"sup_id":  dto.SupID,
	}
	_, err = tx.Exec(ctx, insertRoot, rootArgs)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return errors.Join(err, tx.Rollback(ctx))
	}
	r.log.Log(ctx, core.LevelTrace, "entity insertion succeeded", slog.Any("dto", dto))
	return tx.Commit(ctx)
}

func (r *repoPgx) SelectByID(rid id.ADT) (Snap, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, selectById, rid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return Snap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[snapData])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return Snap{}, err
	}
	r.log.Log(ctx, core.LevelTrace, "entity selection succeeded", slog.Any("dto", dto))
	return DataToSnap(dto)
}

func (r *repoPgx) SelectAll() ([]Ref, error) {
	query := `
		select
			pool_id, rev, title
		from pool_roots`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[refData])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return nil, err
	}
	return DataToRefs(dtos)
}

const (
	selectById = `
		select
			tr.pool_id,
			(array_agg(tr.title))[1] as title,
			jsonb_agg(to_jsonb((select sub from (select ts.pool_id, ts.rev, ts.title) sub))) filter (where ts.pool_id is not null) as subs
		from pool_roots tr
		left join pool_roots ts
			on ts.sup_id = tr.pool_id
		where tr.pool_id = $1
		group by tr.pool_id`
)
