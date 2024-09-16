package state

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

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
				"msg_id":  tr.MsgID,
				"msg_key": tr.MsgKey,
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
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[*RefData])
	if err != nil {
		return nil, err
	}
	return DataToRefs(dtos)
}

func (r *repoPgx) SelectByID(sid id.ADT[ID]) (Root, error) {
	fooId := id.New[ID]()
	queue := &WithRoot{
		ID: id.New[ID](),
		Choices: map[Label]Root{
			"enq": &TensorRoot{
				ID: id.New[ID](),
				S:  &RecRoot{ID: fooId, Name: "Foo"},
				T:  &RecRoot{ID: sid, Name: "Queue"},
			},
			"deq": &PlusRoot{
				ID: id.New[ID](),
				Choices: map[Label]Root{
					"some": &LolliRoot{
						ID: id.New[ID](),
						S:  &RecRoot{ID: fooId, Name: "Foo"},
						T:  &RecRoot{ID: sid, Name: "Queue"},
					},
					"none": &OneRoot{ID: id.New[ID]()},
				},
			},
		},
	}
	return queue, nil
}
