package chnl

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
	name := slog.String("name", "chnlRepoPgx")
	return &repoPgx{p, l.With(name)}
}

func (r *repoPgx) Insert(root Root) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto, err := DataFromRoot(root)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO channels (
			id, name, pre_id, pak, cak, state
		) VALUES (
			@id, @name, @pre_id, @pak, @cak, @state
		)`
	args := pgx.NamedArgs{
		"id":     dto.ID,
		"name":   dto.Name,
		"pre_id": dto.PreID,
		"pak":    dto.PAK,
		"cak":    dto.CAK,
		"state":  dto.St,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *repoPgx) InsertCtx(roots []Root) (rs []Root, err error) {
	query := `
		INSERT INTO channels (
			id, name, pre_id, pak, cak, state
		)
		SELECT
			@new_id, name, @pre_id, pak, @cak, state
		FROM channels
		WHERE id = @id
		RETURNING id, name, pre_id, pak, cak, state`
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	batch := pgx.Batch{}
	dtos, err := DataFromRoots(roots)
	if err != nil {
		return nil, err
	}
	for _, dto := range dtos {
		args := pgx.NamedArgs{
			"new_id": dto.ID,
			"pre_id": dto.PreID,
			"cak":    dto.CAK,
			"id":     dto.PreID,
		}
		batch.Queue(query, args)
	}
	br := tx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var foes []rootData
	for _, dto := range dtos {
		// var foo rootData
		// _, err = br.Exec()
		// err = br.QueryRow().Scan(&foo)
		rows, err := br.Query()
		if err != nil {
			r.log.Error("query execution failed",
				slog.Any("reason", err),
				slog.Any("dto", dto),
			)
		}
		foo, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
		if err != nil {
			r.log.Error("row collection failed",
				slog.Any("reason", err),
				slog.Any("dto", dto),
			)
		}
		foes = append(foes, foo)
	}
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	err = br.Close()
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}
	r.log.Debug("context insertion succeeded", slog.Any("dtos", dtos))
	r.log.Log(ctx, core.LevelTrace, "context insertion succeeded", slog.Any("dtos", dtos))
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	return DataToRoots(foes)
}

func (r *repoPgx) SelectAll() ([]Ref, error) {
	roots := make([]Ref, 5)
	return roots, nil
}

func (r *repoPgx) SelectByID(rid id.ADT) (Root, error) {
	query := `
		SELECT
			id, name, pre_id, pak, cak, state
		FROM channels
		WHERE id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
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
	r.log.Log(ctx, core.LevelTrace, "channel selection succeeded", slog.Any("dto", dto))
	return DataToRoot(dto)
}
