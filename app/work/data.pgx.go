package work

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"
)

// Adapter
type workRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newWorkRepoPgx(p *pgxpool.Pool, l *slog.Logger) *workRepoPgx {
	name := slog.String("name", "workRepoPgx")
	return &workRepoPgx{p, l.With(name)}
}

func (r *workRepoPgx) Insert(root WorkRoot) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := DataFromWorkRoot(root)
	query := `
		INSERT INTO works (
			id, name
		) VALUES (
			@id, @name
		)`
	args := pgx.NamedArgs{
		"id":   dto.ID,
		"name": dto.Name,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("work", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *workRepoPgx) SelectById(id core.ID[Work]) (WorkRoot, error) {
	return WorkRoot{ID: id, Name: "WorkRoot"}, nil
}

func (r *workRepoPgx) SelectChildren(id core.ID[Work]) ([]WorkTeaser, error) {
	query := `
		SELECT
			w.id,
			w.name
		FROM works w
		LEFT JOIN kinships k
			ON w.id = k.child_id
		WHERE k.parent_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[workTeaserData])
	if err != nil {
		return nil, err
	}
	return DataToWorkTeasers(dtos)
}

func (r *workRepoPgx) SelectAll() ([]WorkTeaser, error) {
	roots := make([]WorkTeaser, 5)
	for i := range 5 {
		roots[i] = WorkTeaser{ID: core.New[Work](), Name: fmt.Sprintf("WorkRoot%v", i)}
	}
	return roots, nil
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
