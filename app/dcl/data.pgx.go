package dcl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"
)

// Adapter
type tpRepoPgx struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newTpRepoPgx(p *pgxpool.Pool, l *slog.Logger) *tpRepoPgx {
	name := slog.String("name", "dcl.tpRepoPgx")
	return &tpRepoPgx{p, l.With(name)}
}

func (r *tpRepoPgx) Insert(tp TpRoot) (err error) {
	ctx := context.Background()
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	data := dataFromTpRoot(tp)
	// states
	sq := `
		INSERT INTO states (
			kind, id, name
		) VALUES (
			@kind, @id, @name
		)`
	sb := pgx.Batch{}
	for _, s := range data.States {
		sa := pgx.NamedArgs{
			"kind": s.K,
			"id":   s.ID,
			"name": s.Name,
		}
		sb.Queue(sq, sa)
	}
	sbr := tx.SendBatch(ctx, &sb)
	defer func() {
		err = errors.Join(err, sbr.Close())
	}()
	for _, s := range data.States {
		_, err = sbr.Exec()
		if err != nil {
			r.log.Error("insert failed", slog.Any("reason", err), slog.Any("state", s))
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
	for _, trs := range data.Trs {
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
	for _, trs := range data.Trs {
		for _, tr := range trs {
			_, err = tbr.Exec()
			if err != nil {
				r.log.Error("insert failed", slog.Any("reason", err), slog.Any("tr", tr))
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
	// root
	rq := `
		INSERT INTO tps (
			id, name
		) VALUES (
			@id, @name
		)`
	ra := pgx.NamedArgs{
		"id":   data.ID,
		"name": data.Name,
	}
	_, err = tx.Exec(ctx, rq, ra)
	if err != nil {
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("tp", ra))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *tpRepoPgx) SelectById(id core.ID[AR]) (TpRoot, error) {
	fooId := core.New[AR]()
	queue := With{
		ID: core.New[AR](),
		Chs: Choices{
			"enq": Tensor{
				ID: core.New[AR](),
				S:  TpRef{fooId, "Foo"},
				T:  TpRef{id, "Queue"},
			},
			"deq": Plus{
				ID: core.New[AR](),
				Chs: Choices{
					"some": Lolli{
						ID: core.New[AR](),
						S:  TpRef{fooId, "Foo"},
						T:  TpRef{id, "Queue"},
					},
					"none": One{ID: core.New[AR]()},
				},
			},
		},
	}
	return TpRoot{id, "Queue", queue}, nil
}

func (r *tpRepoPgx) SelectAll() ([]TpRoot, error) {
	query := `
		SELECT
			id,
			name
		FROM tps`
	ctx := context.Background()
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tps, err := pgx.CollectRows(rows, pgx.RowToStructByName[tpRootData])
	if err != nil {
		return nil, err
	}
	return DataToTpRoots(tps)
}

// Adapter
type expRepoPgx struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func newExpRepoPgx(p *pgxpool.Pool, l *slog.Logger) *expRepoPgx {
	name := slog.String("name", "dcl.expRepoPgx")
	return &expRepoPgx{p, l.With(name)}
}

func (r *expRepoPgx) Insert(exp ExpRoot) error {
	return nil
}

func (r *expRepoPgx) SelectById(id core.ID[AR]) (ExpRoot, error) {
	return ExpRoot{ID: id, Name: "ExpRoot"}, nil
}

func (r *expRepoPgx) SelectAll() ([]ExpRoot, error) {
	roots := make([]ExpRoot, 5)
	for i := range 5 {
		roots[i] = ExpRoot{ID: core.New[AR](), Name: fmt.Sprintf("ExpRoot%v", i)}
	}
	return roots, nil
}
