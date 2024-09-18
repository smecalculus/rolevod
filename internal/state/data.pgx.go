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
	dto := dataFromRoot(root)
	// states
	sq := `
		INSERT INTO states (
			id, kind
		) VALUES (
			@id, @kind
		)`
	sb := pgx.Batch{}
	for _, s := range dto.States {
		sa := pgx.NamedArgs{
			"id":   s.ID,
			"kind": s.K,
		}
		sb.Queue(sq, sa)
	}
	sbr := tx.SendBatch(ctx, &sb)
	defer func() {
		err = errors.Join(err, sbr.Close())
	}()
	for _, s := range dto.States {
		_, err = sbr.Exec()
		if err != nil {
			r.log.Error("query execution failed", slog.Any("reason", err), slog.Any("state", s))
		}
	}
	if err != nil {
		return errors.Join(err, sbr.Close(), tx.Rollback(ctx))
	}
	err = sbr.Close()
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	// transitions
	tq := `
		INSERT INTO transitions (
			from_id, to_id, msg_id, msg_key
		) VALUES (
			@from_id, @to_id, @msg_id, @msg_key
		)`
	tb := pgx.Batch{}
	for _, trs := range dto.Trs {
		for _, tr := range trs {
			ta := pgx.NamedArgs{
				"from_id": tr.FromID,
				"to_id":   tr.ToID,
				"msg_id":  tr.OnID,
				"msg_key": tr.OnKey,
			}
			tb.Queue(tq, ta)
		}
	}
	tbr := tx.SendBatch(ctx, &tb)
	defer func() {
		err = errors.Join(err, tbr.Close())
	}()
	for _, trs := range dto.Trs {
		for _, tr := range trs {
			_, err = tbr.Exec()
			if err != nil {
				r.log.Error("query execution failed", slog.Any("reason", err), slog.Any("tr", tr))
			}
		}
	}
	if err != nil {
		return errors.Join(err, tbr.Close(), tx.Rollback(ctx))
	}
	err = tbr.Close()
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
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[*refData])
	if err != nil {
		return nil, err
	}
	return DataToRefs(dtos)
}

func (r *repoPgx) SelectByID(rid id.ADT[ID]) (Root, error) {
	query := `
		SELECT
			s1.id as from_id,
			s1.kind as from_kind,
			tr.msg_id,
			tr.msg_key,
			s2.id as to_id,
			s2.kind as to_kind
		FROM states s1
			LEFT JOIN transitions tr
			ON s1.id = tr.from_id
			LEFT JOIN states s2
			ON tr.to_id = s2.id
		WHERE s1.id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[transition2])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return nil, err
	}
	if len(dtos) == 0 {
		return nil, fmt.Errorf("no rows selected")
	}
	r.log.Log(ctx, core.LevelTrace, "state selection succeeded", slog.Any("dtos", dtos))
	return dataToRoot2(dtos, rid.String()), nil

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
