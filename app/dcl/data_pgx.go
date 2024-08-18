package dcl

import (
	"fmt"
	"log/slog"

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

func (r *tpRepoPgx) Insert(tp TpRoot) error {
	return nil
}

func (r *tpRepoPgx) SelectById(id core.ID[AR]) (TpRoot, error) {
	fooId := core.New[AR]()
	queue := With{
		Choices{
			"enq": Tensor{
				TpName{fooId, "Foo"},
				TpName{id, "Queue"},
			},
			"deq": Plus{
				Choices{
					"some": Lolli{
						TpName{fooId, "Foo"},
						TpName{id, "Queue"},
					},
					"none": One{},
				},
			},
		},
	}
	return TpRoot{id, "Queue", queue}, nil
}

func (r *tpRepoPgx) SelectAll() ([]TpRoot, error) {
	tpDefs := make([]TpRoot, 5)
	for i := range 5 {
		tpDefs[i] = TpRoot{core.New[AR](), fmt.Sprintf("TpRoot%v", i), One{}}
	}
	return tpDefs, nil
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
	expDecs := make([]ExpRoot, 5)
	for i := range 5 {
		expDecs[i] = ExpRoot{core.New[AR](), fmt.Sprintf("ExpRoot%v", i)}
	}
	return expDecs, nil
}
