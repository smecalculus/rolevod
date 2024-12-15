package sig

import (
	"context"
	"errors"
	"log/slog"
	"math"

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
	name := slog.String("name", "sigRepoPgx")
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
	insertRoot := `
		insert into sig_roots (
			sig_id, rev, title
		) VALUES (
			@sig_id, @rev, @title
		)`
	rootArgs := pgx.NamedArgs{
		"sig_id": dto.ID,
		"rev":    dto.Rev,
		"title":  dto.Title,
	}
	_, err = tx.Exec(ctx, insertRoot, rootArgs)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return errors.Join(err, tx.Rollback(ctx))
	}
	insertPE := `
		insert into sig_pes (
			sig_id, rev_from, rev_to, chnl_key, role_fqn
		) VALUES (
			@sig_id, @rev_from, @rev_to, @chnl_key, @role_fqn
		)`
	peArgs := pgx.NamedArgs{
		"sig_id":   dto.ID,
		"rev_from": dto.Rev,
		"rev_to":   math.MaxInt64,
		"chnl_key": dto.PE.Key,
		"role_fqn": dto.PE.Link,
	}
	_, err = tx.Exec(ctx, insertPE, peArgs)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return errors.Join(err, tx.Rollback(ctx))
	}
	insertCE := `
		insert into sig_ces (
			sig_id, rev_from, rev_to, chnl_key, role_fqn
		) VALUES (
			@sig_id, @rev_from, @rev_to, @chnl_key, @role_fqn
		)`
	batch := pgx.Batch{}
	for _, ce := range dto.CEs {
		args := pgx.NamedArgs{
			"sig_id":   dto.ID,
			"rev_from": dto.Rev,
			"rev_to":   math.MaxInt64,
			"chnl_key": ce.Key,
			"role_fqn": ce.Link,
		}
		batch.Queue(insertCE, args)
	}
	br := tx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, ce := range dto.CEs {
		_, err = br.Exec()
		if err != nil {
			r.log.Error("query execution failed", slog.Any("reason", err), slog.Any("ce", ce))
		}
	}
	if err != nil {
		return errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	err = br.Close()
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *repoPgx) SelectByID(rid id.ADT) (Root, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, selectById, rid.String())
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
	r.log.Log(ctx, core.LevelTrace, "signature selection succeeded", slog.Any("dto", dto))
	return DataToRoot(dto)
}

func (r *repoPgx) SelectEnv(ids []ID) (map[ID]Root, error) {
	sigs, err := r.SelectByIDs(ids)
	if err != nil {
		return nil, err
	}
	env := make(map[ID]Root, len(sigs))
	for _, s := range sigs {
		env[s.ID] = s
	}
	return env, nil
}

func (r *repoPgx) SelectByIDs(ids []ID) ([]Root, error) {
	if len(ids) == 0 {
		return []Root{}, nil
	}
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
		batch.Queue(selectById, rid.String())
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
	r.log.Log(ctx, core.LevelTrace, "signatures selection succeeded", slog.Any("dtos", dtos))
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	return DataToRoots(dtos)
}

func (r *repoPgx) SelectChildren(id ID) ([]Ref, error) {
	query := `
		SELECT
			s.id,
			s.name
		FROM signatures s
		LEFT JOIN kinships k
			ON s.id = k.child_id
		WHERE k.parent_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, id.String())
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

func (r *repoPgx) SelectAll() ([]Ref, error) {
	query := `
		select
			sig_id, rev, title
		from sig_roots`
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
			sr.sig_id,
			sr.rev,
			(array_agg(sr.title))[1] as title,
			(jsonb_agg(to_jsonb((select ep from (select sp.chnl_key, sp.role_fqn) ep))))[0] as pe,
			jsonb_agg(to_jsonb((select ep from (select sc.chnl_key, sc.role_fqn) ep))) filter (where sc.sig_id is not null) as ces
		from sig_roots sr
		left join sig_pes sp
			on sp.sig_id = sr.sig_id
			and sp.rev_from >= sr.rev
			and sp.rev_to > sr.rev
		left join sig_ces sc
			on sc.sig_id = sr.sig_id
			and sc.rev_from >= sr.rev
			and sc.rev_to > sr.rev
		where sr.sig_id = $1
		group by sr.sig_id, sr.rev`
)
