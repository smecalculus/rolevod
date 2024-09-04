package role

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/core"
)

// Adapter
type roleRepoPgx struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func newRoleRepoPgx(p *pgxpool.Pool, l *slog.Logger) *roleRepoPgx {
	name := slog.String("name", "roleRepoPgx")
	return &roleRepoPgx{p, l.With(name)}
}

func (r *roleRepoPgx) Insert(rr RoleRoot) (err error) {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto := dataFromRoleRoot(rr)
	// states
	sq := `
		INSERT INTO states (
			kind, id, name
		) VALUES (
			@kind, @id, @name
		)`
	sb := pgx.Batch{}
	for _, s := range dto.States {
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
	for _, s := range dto.States {
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
	// role
	rq := `
		INSERT INTO roles (
			id, name
		) VALUES (
			@id, @name
		)`
	ra := pgx.NamedArgs{
		"id":   dto.ID,
		"name": dto.Name,
	}
	_, err = tx.Exec(ctx, rq, ra)
	if err != nil {
		r.log.Error("insert failed", slog.Any("reason", err), slog.Any("role", ra))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *roleRepoPgx) SelectById(id core.ID[Role]) (RoleRoot, error) {
	fooId := core.New[Role]()
	queue := With{
		ID: core.New[Role](),
		Chs: Choices{
			"enq": Tensor{
				ID: core.New[Role](),
				S:  TpRef{fooId, "Foo"},
				T:  TpRef{id, "Queue"},
			},
			"deq": Plus{
				ID: core.New[Role](),
				Chs: Choices{
					"some": Lolli{
						ID: core.New[Role](),
						S:  TpRef{fooId, "Foo"},
						T:  TpRef{id, "Queue"},
					},
					"none": One{ID: core.New[Role]()},
				},
			},
		},
	}
	return RoleRoot{ID: id, Name: "Queue", St: queue}, nil
}

func (r *roleRepoPgx) SelectChildren(id core.ID[Role]) ([]RoleTeaser, error) {
	query := `
		SELECT
			r.id,
			r.name
		FROM roles r
		LEFT JOIN kinships k
			ON r.id = k.child_id
		WHERE k.parent_id = $1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[roleTeaserData])
	if err != nil {
		return nil, err
	}
	return DataToRoleTeasers(dtos)
}

func (r *roleRepoPgx) SelectAll() ([]RoleTeaser, error) {
	query := `
		SELECT
			id,
			name
		FROM roles`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[roleTeaserData])
	if err != nil {
		return nil, err
	}
	return DataToRoleTeasers(dtos)
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

func (r *kinshipRepoPgx) Insert(kr KinshipRoot) error {
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
	dto := DataFromKinshipRoot(kr)
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
