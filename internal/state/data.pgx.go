package state

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
	name := slog.String("name", "stateRepoPgx")
	return &repoPgx{p, l.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return &repoPgx{}
}

func (r *repoPgx) Insert(root Root) (err error) {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := dataFromRoot2(root)
	query := `
		INSERT INTO states (
			id, kind, from_id, fqn, pair, choices, type
		) VALUES (
			@id, @kind, @from_id, @fqn, @pair, @choices, @type
		)`
	batch := pgx.Batch{}
	for _, st := range dto.States {
		sa := pgx.NamedArgs{
			"id":      st.ID,
			"kind":    st.K,
			"fqn":     st.FQN,
			"from_id": st.FromID,
			"pair":    st.Pair,
			"choices": st.Choices,
			"type":    st.Type,
		}
		batch.Queue(query, sa)
	}
	br := tx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, st := range dto.States {
		_, err = br.Exec()
		if err != nil {
			r.log.Error("query execution failed", slog.Any("reason", err), slog.Any("state", st))
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

func (r *repoPgx) SelectAll() ([]Ref, error) {
	query := `
		SELECT
			kind, id
		FROM states`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[*RefData])
	if err != nil {
		return nil, err
	}
	return DataToRefs(dtos)
}

func (r *repoPgx) SelectByID(rid id.ADT) (Root, error) {
	query := `
		WITH RECURSIVE top_states AS (
			SELECT
				rs.*
			FROM states rs
			WHERE id = $1
			UNION ALL
			SELECT
				bs.*
			FROM states bs, top_states ts
			WHERE bs.from_id = ts.id
		)
		SELECT * FROM top_states`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[state])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return nil, err
	}
	if len(dtos) == 0 {
		err := fmt.Errorf("no rows selected")
		r.log.Error("state selection failed", slog.Any("reason", err))
		return nil, err
	}
	r.log.Log(ctx, core.LevelTrace, "state selection succeeded", slog.Any("dtos", dtos))
	states := make(map[string]state, len(dtos))
	for _, dto := range dtos {
		states[dto.ID] = dto
	}
	return statesToRoot2(states, states[rid.String()])
}

func (r *repoPgx) SelectEnv(ids []ID) (map[ID]Root, error) {
	states, err := r.SelectByIDs(ids)
	if err != nil {
		return nil, err
	}
	env := make(map[ID]Root, len(states))
	for _, st := range states {
		env[st.RID()] = st
	}
	return env, nil
}

func (r *repoPgx) SelectByIDs(ids []ID) (rs []Root, err error) {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	batch := pgx.Batch{}
	for _, rid := range ids {
		batch.Queue(selectByID, rid.String())
	}
	br := tx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var roots []Root
	for _, rid := range ids {
		rows, err := br.Query()
		if err != nil {
			r.log.Error("query execution failed",
				slog.Any("reason", err),
				slog.Any("id", rid),
			)
		}
		dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[state])
		if err != nil {
			r.log.Error("row collection failed",
				slog.Any("reason", err),
				slog.Any("id", rid),
			)
		}
		if len(dtos) == 0 {
			err = ErrDoesNotExist(rid)
			r.log.Error("state selection failed",
				slog.Any("reason", err),
			)
		}
		root, err := dataToRoot(&rootData2{rid.String(), dtos})
		if err != nil {
			r.log.Error("dto mapping failed",
				slog.Any("reason", err),
				slog.Any("dtos", dtos),
			)
		}
		roots = append(roots, root)
	}
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	err = br.Close()
	if err != nil {
		return nil, errors.Join(err, tx.Rollback(ctx))
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Join(err, br.Close(), tx.Rollback(ctx))
	}
	r.log.Log(ctx, core.LevelTrace, "states selection succeeded", slog.Any("roots", roots))
	return roots, err
}

const (
	selectByID = `
		WITH RECURSIVE state_tree AS (
			SELECT root.*
			FROM states root
			WHERE id = $1
			UNION ALL
			SELECT child.*
			FROM states child, state_tree parent
			WHERE child.from_id = parent.id
		)
		SELECT * FROM state_tree
	`
)
