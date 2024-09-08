package agent

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
type agentRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newAgentRepoPgx(p *pgxpool.Pool, l *slog.Logger) *agentRepoPgx {
	name := slog.String("name", "agentRepoPgx")
	return &agentRepoPgx{p, l.With(name)}
}

func (r *agentRepoPgx) Insert(root AgentRoot) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := DataFromAgentRoot(root)
	query := `
		INSERT INTO agents (
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
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("agent", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *agentRepoPgx) SelectById(id id.ADT[ID]) (AgentRoot, error) {
	return AgentRoot{ID: id, Name: "AgentRoot"}, nil
}

func (r *agentRepoPgx) SelectChildren(id id.ADT[ID]) ([]AgentRef, error) {
	query := `
		SELECT
			f.id,
			f.name
		FROM agents f
		LEFT JOIN kinships k
			ON f.id = k.child_id
		WHERE k.parent_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[agentRefData])
	if err != nil {
		return nil, err
	}
	return DataToAgentRefs(dtos)
}

func (r *agentRepoPgx) SelectAll() ([]AgentRef, error) {
	roots := make([]AgentRef, 5)
	for i := range 5 {
		roots[i] = AgentRef{ID: id.New[ID](), Name: fmt.Sprintf("AgentRoot%v", i)}
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
