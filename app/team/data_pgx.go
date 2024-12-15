package team

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
	name := slog.String("name", "teamRepoPgx")
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
		INSERT INTO team_roots (
			team_id, rev, title, sup_id
		) VALUES (
			@team_id, @rev, @title, @sup_id
		)`
	rootArgs := pgx.NamedArgs{
		"team_id": dto.ID,
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
			team_id, rev, title
		from team_roots`
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
			tr.team_id,
			tr.rev,
			(array_agg(tr.title))[1] as title,
			jsonb_agg(to_jsonb((select sub from (select ts.team_id, ts.rev, ts.title) sub))) filter (where ts.team_id is not null) as subs
		from team_roots tr
		left join team_roots ts
			on ts.sup_id = tr.team_id
		where tr.team_id = $1
		group by tr.team_id, tr.rev`
)
