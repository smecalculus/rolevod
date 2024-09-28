package deal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/app/seat"
)

// Adapter
type dealRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newDealRepoPgx(p *pgxpool.Pool, l *slog.Logger) *dealRepoPgx {
	name := slog.String("name", "dealRepoPgx")
	return &dealRepoPgx{p, l.With(name)}
}

func (r *dealRepoPgx) Insert(root DealRoot) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := DataFromDealRoot(root)
	query := `
		INSERT INTO deals (
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
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("deal", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *dealRepoPgx) SelectAll() ([]DealRef, error) {
	roots := make([]DealRef, 5)
	for i := range 5 {
		roots[i] = DealRef{ID: id.New(), Name: fmt.Sprintf("DealRoot%v", i)}
	}
	return roots, nil
}

func (r *dealRepoPgx) SelectByID(id id.ADT) (DealRoot, error) {
	return DealRoot{ID: id, Name: "DealRoot"}, nil
}

func (r *dealRepoPgx) SelectChildren(id id.ADT) ([]DealRef, error) {
	query := `
		SELECT
			d.id,
			d.name
		FROM deals d
		LEFT JOIN kinships k
			ON d.id = k.child_id
		WHERE k.parent_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[dealRefData])
	if err != nil {
		return nil, err
	}
	return DataToDealRefs(dtos)
}

func (r *dealRepoPgx) SelectSeats(id id.ADT) ([]seat.SeatRef, error) {
	return []seat.SeatRef{}, nil
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

// Adapter
type partRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newPartRepoPgx(p *pgxpool.Pool, l *slog.Logger) *partRepoPgx {
	name := slog.String("name", "partRepoPgx")
	return &partRepoPgx{p, l.With(name)}
}

func (r *partRepoPgx) Insert(root PartRoot) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := DataFromPartRoot(root)
	query := `
		INSERT INTO participations (
			part_id, deal_id, seat_id, pak, cak
		) values (
			@part_id, @deal_id, @seat_id, @pak, @cak
		)`
	args := pgx.NamedArgs{
		"part_id": dto.PartID,
		"deal_id": dto.DealID,
		"seat_id": dto.SeatID,
		"pak":     dto.PAK,
		"cak":     dto.CAK,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("participation insertion failed",
			slog.Any("reason", err),
			slog.Any("part", args),
		)
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *partRepoPgx) SelectByID(rid ID) (PartRoot, error) {
	query := `
		SELECT
			part_id, deal_id, seat_id, pak, cak
		FROM participations
		WHERE part_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return PartRoot{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[partRootData])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return PartRoot{}, err
	}
	r.log.Log(ctx, core.LevelTrace, "participation selection succeeded", slog.Any("dto", dto))
	return DataToPartRoot(dto)
}
