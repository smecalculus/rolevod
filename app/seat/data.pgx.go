package seat

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/id"
)

// Adapter
type seatRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newSeatRepoPgx(p *pgxpool.Pool, l *slog.Logger) *seatRepoPgx {
	name := slog.String("name", "seatRepoPgx")
	return &seatRepoPgx{p, l.With(name)}
}

func (r *seatRepoPgx) Insert(root SeatRoot) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto, err := DataFromSeatRoot(root)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO seats (
			id, name, via, ctx
		) VALUES (
			@id, @name, @via, @ctx
		)`
	args := pgx.NamedArgs{
		"id":   dto.ID,
		"name": dto.Name,
		"via":  dto.Via,
		"ctx":  dto.Ctx,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *seatRepoPgx) SelectByID(id id.ADT[ID]) (SeatRoot, error) {
	query := `
	SELECT
		id, name, via, ctx
	FROM seats
	WHERE id=$1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, id.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return SeatRoot{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[seatRootData])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return SeatRoot{}, err
	}
	r.log.Debug("seat selection succeed", slog.Any("dto", dto))
	return DataToSeatRoot(dto)
}

func (r *seatRepoPgx) SelectChildren(id id.ADT[ID]) ([]SeatRef, error) {
	query := `
		SELECT
			s.id,
			s.name
		FROM seats s
		LEFT JOIN kinships k
			ON s.id = k.child_id
		WHERE k.parent_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[SeatRefData])
	if err != nil {
		return nil, err
	}
	return DataToSeatRefs(dtos)
}

func (r *seatRepoPgx) SelectAll() ([]SeatRef, error) {
	roots := make([]SeatRef, 5)
	for i := range 5 {
		roots[i] = SeatRef{ID: id.New[ID](), Name: fmt.Sprintf("SeatRoot%v", i)}
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
