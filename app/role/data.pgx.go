package role

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"smecalculus/rolevod/lib/id"
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

func (r *roleRepoPgx) Insert(root RoleRoot) error {
	ctx := context.Background()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	dto, err := dataFromRoleRoot(root)
	if err != nil {
		return err
	}
	rq := `
		INSERT INTO roles (
			id, name, state
		) VALUES (
			@id, @name, @state
		)`
	ra := pgx.NamedArgs{
		"id":    dto.ID,
		"name":  dto.Name,
		"state": dto.State,
	}
	_, err = tx.Exec(ctx, rq, ra)
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return errors.Join(err, tx.Rollback(ctx))
	}
	return tx.Commit(ctx)
}

func (r *roleRepoPgx) SelectAll() ([]RoleRef, error) {
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
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[RoleRefData])
	if err != nil {
		return nil, err
	}
	return DataToRoleRefs(dtos)
}

func (r *roleRepoPgx) SelectByID(rid id.ADT[ID]) (RoleRoot, error) {
	// fooID := id.New[state.ID]()
	// queueID, _ := id.String[state.ID](rid.String())
	// queue := &state.WithRoot{
	// 	ID: id.New[state.ID](),
	// 	Choices: map[state.Label]state.Root{
	// 		"enq": &state.TensorRoot{
	// 			ID: id.New[state.ID](),
	// 			S:  &state.TpRefRoot{ID: fooID, Name: "Foo"},
	// 			T:  &state.TpRefRoot{ID: queueID, Name: "Queue"},
	// 		},
	// 		"deq": &state.PlusRoot{
	// 			ID: id.New[state.ID](),
	// 			Choices: map[state.Label]state.Root{
	// 				"some": &state.LolliRoot{
	// 					ID: id.New[state.ID](),
	// 					S:  &state.TpRefRoot{ID: fooID, Name: "Foo"},
	// 					T:  &state.TpRefRoot{ID: queueID, Name: "Queue"},
	// 				},
	// 				"none": &state.OneRoot{ID: id.New[state.ID]()},
	// 			},
	// 		},
	// 	},
	// }
	// return RoleRoot{ID: rid, Name: "Queue", State: state.ConvertRootToRef(queue)}, nil
	query := `
		SELECT
			id, name, state
		FROM roles
		WHERE id=$1`
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, query, rid.String())
	if err != nil {
		r.log.Error("query execution failed", slog.Any("reason", err))
		return RoleRoot{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[roleRootData])
	if err != nil {
		r.log.Error("row collection failed", slog.Any("reason", err))
		return RoleRoot{}, err
	}
	r.log.Debug("role selection succeeded", slog.Any("dto", dto))
	return dataToRoleRoot(dto)
}

func (r *roleRepoPgx) SelectChildren(id id.ADT[ID]) ([]RoleRef, error) {
	query := `
		SELECT
			r.id,
			r.name,
			r.state
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
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[RoleRefData])
	if err != nil {
		return nil, err
	}
	return DataToRoleRefs(dtos)
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
	dto, err := DataFromKinshipRoot(root)
	if err != nil {
		return err
	}
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
