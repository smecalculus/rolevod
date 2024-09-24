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
func (r *repoPgx) Insert(root Root) (err error) {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := dataFromRoot3(root)
	query := `
		INSERT INTO states2 (
			id, kind, from_id, on_ref, to_id, to_ids
		) VALUES (
			@id, @kind, @from_id, @on_ref, @to_id, @to_ids
		)`
	batch := pgx.Batch{}
	for _, s := range dto.States {
		sa := pgx.NamedArgs{
			"id":      s.ID,
			"kind":    s.K,
			"from_id": s.FromID,
			"on_ref":  s.OnRef,
			"to_id":   s.ToID,
			"to_ids":  s.ToIDs,
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
			kind,
			id
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

func (r *repoPgx) SelectByID(rid id.ADT[ID]) (Root, error) {
	query := `
		WITH RECURSIVE sts AS (
			SELECT
				id, kind, from_id, on_ref, to_id, to_ids
			FROM states2
			WHERE id = $1
			UNION ALL
			SELECT
				s1.id, s1.kind, s1.from_id, s1.on_ref, s1.to_id, s1.to_ids
			FROM states2 s1, sts s2
			WHERE s1.from_id = s2.id
		)
		SELECT * FROM sts`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[state2])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return nil, err
	}
	if len(dtos) == 0 {
		return nil, fmt.Errorf("no rows selected")
	}
	r.log.Log(ctx, core.LevelTrace, "state selection succeeded", slog.Any("dtos", dtos))
	return dataToRoot3(dtos, rid.String())

	// fooId := id.New[ID]()
	// queue := &WithRoot{
	// 	ID: id.New[ID](),
	// 	Choices: map[Label]Root{
	// 		"enq": &TensorRoot{
	// 			ID: id.New[ID](),
	// 			A:  &RecurRoot{ID: fooId, Name: "Foo"},
	// 			C:  &RecurRoot{ID: sid, Name: "Queue"},
	// 		},
	// 		"deq": &PlusRoot{
	// 			ID: id.New[ID](),
	// 			Choices: map[Label]Root{
	// 				"some": &LolliRoot{
	// 					ID: id.New[ID](),
	// 					X:  &RecurRoot{ID: fooId, Name: "Foo"},
	// 					Z:  &RecurRoot{ID: sid, Name: "Queue"},
	// 				},
	// 				"none": &OneRoot{ID: id.New[ID]()},
	// 			},
	// 		},
	// 	},
	// }
	// return queue, nil
}
