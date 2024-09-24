package chnl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/ak"
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

func (r *repoPgx) SelectAll() ([]Ref, error) {
	roots := make([]Ref, 5)
	for i := range 5 {
		roots[i] = Ref{
			Name: fmt.Sprintf("Root%v", i),
			PAK:  ak.New(),
			CAK:  ak.New(),
		}
	}
	return roots, nil
}

func (r *repoPgx) SelectByID(rid id.ADT[ID]) (Root, error) {
	query := `
		SELECT
			id, name, pre_id, pak, cak, state
		FROM channels
		WHERE id=$1`
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
