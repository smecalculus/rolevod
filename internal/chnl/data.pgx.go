package chnl

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"
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
	dto, err := DataFromRoot(root)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO channels (
			id, name, pre_id, st_id, state
		) VALUES (
			@id, @name, @pre_id, @st_id, @state
		)`
	args := pgx.NamedArgs{
		"id":     dto.ID,
		"name":   dto.Name,
		"pre_id": dto.PreID,
		"st_id":  dto.StID,
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
			id, name, pre_id, st_id, state
		)
		SELECT
			@new_id, name, @pre_id, st_id, state
		FROM channels
		WHERE id = @id
		RETURNING *`
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	batch := pgx.Batch{}
	reqs, err := DataFromRoots(roots)
	if err != nil {
		r.log.Error("dtos mapping failed",
			slog.Any("reason", err),
			slog.Any("roots", roots),
		)
		return nil, err
	}
	for _, req := range reqs {
		args := pgx.NamedArgs{
			"new_id": req.ID,
			"pre_id": req.PreID,
			"id":     req.PreID,
		}
		batch.Queue(query, args)
	}
	br := tx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var resps []rootData
	for _, req := range reqs {
		rows, err := br.Query()
		if err != nil {
			r.log.Error("query execution failed",
				slog.Any("reason", err),
				slog.Any("req", req),
			)
		}
		resp, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
		if err != nil {
			r.log.Error("row collection failed",
				slog.Any("reason", err),
				slog.Any("req", req),
			)
		}
		resps = append(resps, resp)
	}
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	err = br.Close()
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}
	r.log.Log(ctx, core.LevelTrace, "ctx insertion succeeded", slog.Any("resps", resps))
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	return DataToRoots(resps)
}

func (r *repoPgx) SelectAll() ([]Ref, error) {
	roots := make([]Ref, 5)
	return roots, nil
}

func (r *repoPgx) SelectByID(rid ID) (Root, error) {
	query := `
		SELECT
			id, name, pre_id, st_id, state
		FROM channels
		WHERE id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err), slog.Any("id", rid))
		return Root{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err), slog.Any("id", rid))
		return Root{}, err
	}
	r.log.Log(ctx, core.LevelTrace, "channel selection succeeded", slog.Any("dto", dto))
	return DataToRoot(dto)
}

func (r *repoPgx) SelectCfg(ids []ID) (map[ID]Root, error) {
	chnls, err := r.SelectMany(ids)
	if err != nil {
		return nil, err
	}
	cfg := make(map[ID]Root, len(chnls))
	for _, ch := range chnls {
		cfg[ch.ID] = ch
	}
	return cfg, nil
}

func (r *repoPgx) SelectMany(ids []ID) (rs []Root, err error) {
	query := `
		SELECT
			id, name, pre_id, st_id, state
		FROM channels
		WHERE id = $1`
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	batch := pgx.Batch{}
	for _, rid := range ids {
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
	r.log.Log(ctx, core.LevelTrace, "channels selection succeeded", slog.Any("dtos", dtos))
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	return DataToRoots(dtos)
}
