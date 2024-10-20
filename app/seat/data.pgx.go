package seat

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
			id, name, pe, ces
		) VALUES (
			@id, @name, @pe, @ces
		)`
	args := pgx.NamedArgs{
		"id":   dto.ID,
		"name": dto.Name,
		"pe":  dto.PE,
		"ces":  dto.CEs,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *seatRepoPgx) SelectByID(rid ID) (SeatRoot, error) {
	query := `
		SELECT
			id, name, pe, ces
		FROM seats
		WHERE id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
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
	r.log.Log(ctx, core.LevelTrace, "seat selection succeeded", slog.Any("dto", dto))
	return DataToSeatRoot(dto)
}

func (r *seatRepoPgx) SelectEnv(ids []ID) (map[ID]SeatRoot, error) {
	seats, err := r.SelectByIDs(ids)
	if err != nil {
		return nil, err
	}
	env := make(map[ID]SeatRoot, len(seats))
	for _, s := range seats {
		env[s.ID] = s
	}
	return env, nil
}

func (r *seatRepoPgx) SelectByIDs(ids []ID) ([]SeatRoot, error) {
	if len(ids) == 0 {
		return []SeatRoot{}, nil
	}
	query := `
		SELECT
			id, name, pe, ces
		FROM seats
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
	var dtos []seatRootData
	for _, rid := range ids {
		rows, err := br.Query()
		if err != nil {
			r.log.Error("query execution failed",
				slog.Any("reason", err),
				slog.Any("id", rid),
			)
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[seatRootData])
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
	r.log.Log(ctx, core.LevelTrace, "seats selection succeeded", slog.Any("dtos", dtos))
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	return DataToSeatRoots(dtos)
}

func (r *seatRepoPgx) SelectChildren(id ID) ([]SeatRef, error) {
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
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[seatRefData])
	if err != nil {
		return nil, err
	}
	return DataToSeatRefs(dtos)
}

func (r *seatRepoPgx) SelectAll() ([]SeatRef, error) {
	return []SeatRef{}, nil
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
