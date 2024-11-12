package role

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
	name := slog.String("name", "roleRepoPgx")
	return &repoPgx{p, l.With(name)}
}

func (r *repoPgx) Insert(root Root) error {
	ctx := context.Background()
	r.log.Log(ctx, core.LevelTrace, "entity insertion started", slog.Any("id", root.ID))
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto, err := DataFromRoot(root)
	if err != nil {
		r.log.Error("dto mapping failed", slog.Any("reason", err))
		return err
	}
	query := `
		INSERT INTO roles (
			id, rev, name, state_id, whole_id
		) VALUES (
			@id, @rev, @name, @state_id, @whole_id
		)`
	args := pgx.NamedArgs{
		"id":       dto.ID,
		"rev":      dto.Rev,
		"name":     dto.Name,
		"state_id": dto.StateID,
		"whole_id": dto.WholeID,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return errors.Join(err, tx.Rollback(ctx))
	}
	r.log.Log(ctx, core.LevelTrace, "entity insertion succeeded", slog.Any("dto", dto))
	return tx.Commit(ctx)
}

func (r *repoPgx) SelectAll() ([]Ref, error) {
	query := `
		SELECT
			id, name
		FROM roles`
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

func (r *repoPgx) SelectByID(eid ID) (Root, error) {
	query := `
		select
			id, rev, name, state_id, whole_id
		from roles
		where id = $1
		order by rev desc
		limit 1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, eid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return Root{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return Root{}, err
	}
	r.log.Log(ctx, core.LevelTrace, "entity selection succeeded", slog.Any("dto", dto))
	return DataToRoot(dto)
}

func (r *repoPgx) SelectEnv(ids []ID) (map[ID]Root, error) {
	roots, err := r.SelectByIDs(ids)
	if err != nil {
		return nil, err
	}
	env := make(map[ID]Root, len(roots))
	for _, root := range roots {
		env[root.ID] = root
	}
	return env, nil
}

func (r *repoPgx) SelectByIDs(ids []ID) ([]Root, error) {
	if len(ids) == 0 {
		return []Root{}, nil
	}
	query := `
		select
			id, rev, name, state_id, whole_id
		from roles
		where id = $1
		order by rev desc
		limit 1`
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	batch := pgx.Batch{}
	for _, rid := range ids {
		if rid.IsEmpty() {
			return nil, id.ErrEmpty
		}
		batch.Queue(query, rid.String())
	}
	br := tx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []rootData
	for _, rid := range ids {
		rows, err := br.Query()
		if err != nil {
			r.log.Error("query execution failed",
				slog.Any("reason", err),
				slog.Any("id", rid),
			)
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
		if err != nil {
			r.log.Error("row collection failed",
				slog.Any("reason", err),
				slog.Any("id", rid),
			)
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	err = br.Close()
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}
	r.log.Log(ctx, core.LevelTrace, "roots selection succeeded", slog.Any("dtos", dtos))
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	return DataToRoots(dtos)
}

func (r *repoPgx) SelectParts(eid id.ADT) ([]Ref, error) {
	query := `
		SELECT
			r.id,
			r.name,
			r.state
		FROM roles r
		LEFT JOIN kinships k
			ON r.id = k.child_id
		WHERE k.parent_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, eid.String())
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
