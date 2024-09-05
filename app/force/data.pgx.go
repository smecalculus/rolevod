package force

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
type forceRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newForceRepoPgx(p *pgxpool.Pool, l *slog.Logger) *forceRepoPgx {
	name := slog.String("name", "forceRepoPgx")
	return &forceRepoPgx{p, l.With(name)}
}

func (r *forceRepoPgx) Insert(root ForceRoot) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := DataFromForceRoot(root)
	query := `
		INSERT INTO forces (
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
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("force", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *forceRepoPgx) SelectById(id core.ID[Force]) (ForceRoot, error) {
	return ForceRoot{ID: id, Name: "ForceRoot"}, nil
}

func (r *forceRepoPgx) SelectChildren(id core.ID[Force]) ([]ForceTeaser, error) {
	query := `
		SELECT
			f.id,
			f.name
		FROM forces f
		LEFT JOIN kinships k
			ON f.id = k.child_id
		WHERE k.parent_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[forceTeaserData])
	if err != nil {
		return nil, err
	}
	return DataToForceTeasers(dtos)
}

func (r *forceRepoPgx) SelectAll() ([]ForceTeaser, error) {
	roots := make([]ForceTeaser, 5)
	for i := range 5 {
		roots[i] = ForceTeaser{ID: core.New[Force](), Name: fmt.Sprintf("ForceRoot%v", i)}
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
