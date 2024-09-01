package ws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"

	"smecalculus/rolevod/app/dcl"
)

// Adapter
type envRepoPgx struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newEnvRepoPgx(p *pgxpool.Pool, l *slog.Logger) *envRepoPgx {
	name := slog.String("name", "ws.envRepoPgx")
	return &envRepoPgx{p, l.With(name)}
}

func (r *envRepoPgx) Insert(root EnvRoot) error {
	query := `
		INSERT INTO envs (
			id,
			name
		) values (
			@id,
			@name
		)`
	data := dataFromEnvRoot(root)
	args := pgx.NamedArgs{
		"id":   data.ID,
		"name": data.Name,
	}
	ctx := context.Background()
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("env", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *envRepoPgx) SelectById(id core.ID[AR]) (EnvRoot, error) {
	query := `
		SELECT
			id,
			name
		FROM envs
		WHERE id = $1`
	ctx := context.Background()
	rows, err := r.conn.Query(ctx, query, core.ToString(id))
	if err != nil {
		return EnvRoot{}, err
	}
	defer rows.Close()
	env, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[envRootData])
	if err != nil {
		return EnvRoot{}, err
	}
	return dataToEnvRoot(env)
}

func (r *envRepoPgx) SelectAll() ([]EnvRoot, error) {
	roots := make([]EnvRoot, 5)
	for i := range 5 {
		roots[i] = EnvRoot{core.New[AR](), fmt.Sprintf("Foo%v", i), []dcl.TpTeaser{}, []dcl.ExpRoot{}}
	}
	return roots, nil
}

// Adapter
type tpRepoPgx struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newTpRepoPgx(p *pgxpool.Pool, l *slog.Logger) *tpRepoPgx {
	name := slog.String("name", "ws.tpRepoPgx")
	return &tpRepoPgx{p, l.With(name)}
}

func (r *tpRepoPgx) Insert(intro TpIntro) error {
	query := `
		INSERT INTO introductions (
			env_id,
			tp_id
		) values (
			@env_id,
			@tp_id
		)`
	data := dataFromTpIntro(intro)
	args := pgx.NamedArgs{
		"env_id": data.EnvID,
		"tp_id":  data.TpID,
	}
	ctx := context.Background()
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("intro", args))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *tpRepoPgx) SelectById(id core.ID[AR]) ([]dcl.TpTeaser, error) {
	query := `
		SELECT
			tp.id,
			tp.name
		FROM tps tp
		LEFT JOIN introductions intro
			ON tp.id = intro.tp_id
		WHERE intro.env_id = $1`
	ctx := context.Background()
	rows, err := r.conn.Query(ctx, query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tps, err := pgx.CollectRows(rows, pgx.RowToStructByName[dcl.TpTeaserData])
	if err != nil {
		return nil, err
	}
	return dcl.DataToTpTeasers(tps)
}
