package sig

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
type sigRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newSigRepoPgx(p *pgxpool.Pool, l *slog.Logger) *sigRepoPgx {
	name := slog.String("name", "sigRepoPgx")
	return &sigRepoPgx{p, l.With(name)}
}

func (r *sigRepoPgx) Insert(root Root) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto, err := DataFromSigRoot(root)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO signatures (
			id, name, pe, ces
		) VALUES (
			@id, @name, @pe, @ces
		)`
	args := pgx.NamedArgs{
		"id":   dto.ID,
		"name": dto.Name,
		"pe":   dto.PE,
		"ces":  dto.CEs,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *sigRepoPgx) SelectByID(rid ID) (Root, error) {
	query := `
		SELECT
			id, name, pe, ces
		FROM signatures
		WHERE id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return Root{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[sigRootData])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return Root{}, err
	}
	r.log.Log(ctx, core.LevelTrace, "signature selection succeeded", slog.Any("dto", dto))
	return DataToSigRoot(dto)
}

func (r *sigRepoPgx) SelectEnv(ids []ID) (map[ID]Root, error) {
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

func (r *sigRepoPgx) SelectByIDs(ids []ID) ([]Root, error) {
	if len(ids) == 0 {
		return []Root{}, nil
	}
	query := `
		SELECT
			id, name, pe, ces
		FROM signatures
		WHERE id = $1`
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
	var dtos []sigRootData
	for _, rid := range ids {
		rows, err := br.Query()
		if err != nil {
			r.log.Error("query execution failed",
				slog.Any("reason", err),
				slog.Any("id", rid),
			)
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[sigRootData])
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
	return DataToSigRoots(dtos)
}

func (r *sigRepoPgx) SelectChildren(id ID) ([]Ref, error) {
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
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[sigRefData])
	if err != nil {
		return nil, err
	}
	return DataToSigRefs(dtos)
}

func (r *sigRepoPgx) SelectAll() ([]Ref, error) {
	return []Ref{}, nil
}

// Adapter
type kinshipRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newKinshipRepoPgx(p *pgxpool.Pool, l *slog.Logger) *kinshipRepoPgx {
	name := slog.String("name", "kinshipRepoPgx")
	return &kinshipRepoPgx{p, l.With(name)}
}

func (r *kinshipRepoPgx) Insert(root KinshipRoot) error {
	query := `
		INSERT INTO kinships (
			parent_id,
			child_id
		) values (
			@parent_id,
			@child_id
		)`
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	batch := pgx.Batch{}
	dto := DataFromKinshipRoot(root)
	for _, child := range dto.Children {
		args := pgx.NamedArgs{
			"parent_id": dto.Parent.ID,
			"child_id":  child.ID,
		}
		batch.Queue(query, args)
	}
	br := tx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, child := range dto.Children {
		_, err = br.Exec()
		if err != nil {
			r.log.Error("insert failed",
				slog.Any("reason", err),
				slog.Any("parent", dto.Parent),
				slog.Any("child", child))
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
