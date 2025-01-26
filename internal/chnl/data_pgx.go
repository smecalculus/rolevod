package chnl

import (
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/id"
)

// Adapter
type repoPgx struct {
	log *slog.Logger
}

func newRepoPgx(l *slog.Logger) *repoPgx {
	name := slog.String("name", "chnlRepoPgx")
	return &repoPgx{l.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return &repoPgx{}
}

func (r *repoPgx) Insert(source data.Source, root Root) error {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", root.ID)
	dto, err := DataFromRoot(root)
	if err != nil {
		return err
	}
	rootArgs := pgx.NamedArgs{
		"chnl_id": dto.ID,
		"title":   dto.Title,
	}
	_, err = ds.Conn.Exec(ds.Ctx, rootInsert, rootArgs)
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", rootInsert))
		return err
	}
	stateArgs := pgx.NamedArgs{
		"chnl_id":  dto.ID,
		"state_id": dto.StateID,
		"rev":      0,
	}
	_, err = ds.Conn.Exec(ds.Ctx, stateInsert, stateArgs)
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", stateInsert))
		return err
	}
	poolArgs := pgx.NamedArgs{
		"chnl_id": dto.ID,
		"pool_id": dto.PoolID,
		"rev":     0,
	}
	_, err = ds.Conn.Exec(ds.Ctx, poolInsert, poolArgs)
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", poolInsert))
		return err
	}
	return nil
}

func (r *repoPgx) UpdateState(source data.Source, root Root) error {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", root.ID)
	dto, err := DataFromRoot(root)
	if err != nil {
		return err
	}
	stateArgs := pgx.NamedArgs{
		"chnl_id":  dto.ID,
		"state_id": dto.StateID,
		"rev":      dto.Revs[stateRev],
	}
	_, err = ds.Conn.Exec(ds.Ctx, stateInsert, stateArgs)
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", stateInsert))
		return err
	}
	rootArgs := pgx.NamedArgs{
		"chnl_id": dto.ID,
		"rev":     dto.Revs[stateRev],
	}
	ct, err := ds.Conn.Exec(ds.Ctx, rootUpdateState, rootArgs)
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", rootUpdateState))
		return err
	}
	if ct.RowsAffected() == 0 {
		r.log.Error("root update failed", idAttr)
		return errOptimisticUpdate(root.Revs[stateRev] - 1)
	}
	return nil
}

func (r *repoPgx) UpdatePool(source data.Source, root Root) error {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", root.ID)
	dto, err := DataFromRoot(root)
	if err != nil {
		return err
	}
	poolArgs := pgx.NamedArgs{
		"chnl_id": dto.ID,
		"pool_id": dto.PoolID,
		"rev":     dto.Revs[stateRev],
	}
	_, err = ds.Conn.Exec(ds.Ctx, poolInsert, poolArgs)
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", poolInsert))
		return err
	}
	rootArgs := pgx.NamedArgs{
		"chnl_id": dto.ID,
		"rev":     dto.Revs[stateRev],
	}
	ct, err := ds.Conn.Exec(ds.Ctx, rootUpdatePool, rootArgs)
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", rootUpdatePool))
		return err
	}
	if ct.RowsAffected() == 0 {
		r.log.Error("root update failed", idAttr)
		return errOptimisticUpdate(root.Revs[stateRev] - 1)
	}
	return nil
}

func (r *repoPgx) SelectRefs(source data.Source) ([]Ref, error) {
	roots := make([]Ref, 5)
	return roots, nil
}

func (r *repoPgx) SelectByID(source data.Source, rid id.ADT) (Root, error) {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", rid)
	rows, err := ds.Conn.Query(ds.Ctx, rootSelectById, rid.String())
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", rootSelectById))
		return Root{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
	if err != nil {
		r.log.Error("row collection failed", idAttr)
		return Root{}, err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "root selection succeeded", slog.Any("dto", dto))
	return DataToRoot(dto)
}

func (r *repoPgx) SelectByIDs(source data.Source, ids []id.ADT) (_ []Root, err error) {
	ds := data.MustConform[data.SourcePgx](source)
	if len(ids) == 0 {
		return []Root{}, nil
	}
	batch := pgx.Batch{}
	for _, rid := range ids {
		if rid.IsEmpty() {
			return nil, id.ErrEmpty
		}
		batch.Queue(rootSelectById, rid.String())
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []rootData
	for _, rid := range ids {
		rows, err := br.Query()
		if err != nil {
			r.log.Error("query execution failed", slog.Any("id", rid), slog.String("q", rootSelectById))
		}
		defer rows.Close()
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rootData])
		if err != nil {
			r.log.Error("row collection failed", slog.Any("id", rid))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "roots selection succeeded", slog.Any("dtos", dtos))
	return DataToRoots(dtos)
}

const (
	rootInsert = `
		insert into chnl_roots (
			chnl_id, title, revs
		) values (
			@chnl_id, @title, '{0, 0}'
		)`
	stateInsert = `
		insert into chnl_states (
			chnl_id, state_id, rev
		) values (
			@chnl_id, @state_id, @rev
		)`
	poolInsert = `
		insert into chnl_pools (
			chnl_id, pool_id, rev
		) values (
			@chnl_id, @pool_id, @rev
		)`
	rootSelectById = `
		select
			chnl_id, title, revs
		from chnl_roots
		where chnl_id = $1`
	snapSelectById = `
		select
			root.chnl_id,
			root.title,
			state.state_id,
			pool.pool_id
		from chnl_roots root
		left join chnl_states state
			on state.chnl_id = root.chnl_id
			and state.rev = root.revs[1]
		left join chnl_pools pool
			on pool.chnl_id = root.chnl_id
			and pool.rev = root.revs[2]
		where root.chnl_id = $1`
	rootUpdateState = `
		update chnl_roots
		set revs[1] = @rev
		where chnl_id = @chnl_id
			and revs[1] = @rev - 1`
	rootUpdatePool = `
		update chnl_roots
		set revs[2] = @rev
		where chnl_id = @chnl_id
			and revs[2] = @rev - 1`
)
