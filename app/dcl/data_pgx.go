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
	data := dataFromTpRoot(tp)
	ctx := context.Background()
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	// states
	sq := "insert into states (kind, id, name) values (@kind, @id, @name)"
	sb := pgx.Batch{}
	for _, s := range data.States {
		sa := pgx.NamedArgs{
			"kind": s.K,
			"id":   s.ID,
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
	tq := "insert into transitions (from_id, to_id, msg_id, label) values (@from, @to, @msg, @label)"
	tb := pgx.Batch{}
	for _, trs := range data.Trs {
		for _, tr := range trs {
			ta := pgx.NamedArgs{
				"from":  tr.FromID,
				"to":    tr.ToID,
				"msg":   tr.MsgID,
				"label": tr.Label,
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
				r.log.Error("insert failed", slog.Any("reason", err), slog.Any("transition", tr))
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
	rq := "insert into tps (id, name) values (@id, @name)"
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
				S:  TpRef{"Foo", fooId},
				T:  TpRef{"Queue", id},
			},
			"deq": Plus{
				ID: core.New[AR](),
				Chs: Choices{
					"some": Lolli{
						ID: core.New[AR](),
						S:  TpRef{"Foo", fooId},
						T:  TpRef{"Queue", id},
					},
					"none": One{ID: core.New[AR]()},
				},
			},
		},
	}
	return TpRoot{id, "Queue", queue}, nil
}

func (r *tpRepoPgx) SelectAll() ([]TpRoot, error) {
	roots := make([]TpRoot, 5)
	for i := range 5 {
		roots[i] = TpRoot{core.New[AR](), fmt.Sprintf("TpRoot%v", i), One{}}
	}
	return roots, nil
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
	return ExpRoot{id, "ExpRoot"}, nil
}

func (r *expRepoPgx) SelectAll() ([]ExpRoot, error) {
	roots := make([]ExpRoot, 5)
	for i := range 5 {
		roots[i] = ExpRoot{core.New[AR](), fmt.Sprintf("ExpRoot%v", i)}
	}
	return roots, nil
}

func (r *tpRepoPgx) WithinTransaction(act func(ctx context.Context) error) error {
	ctx := context.Background()
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	err = act(context.WithValue(ctx, "tx", tx))
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}
