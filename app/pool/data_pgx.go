package pool

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/proc"
	"smecalculus/rolevod/internal/step"
)

// Adapter
type repoPgx struct {
	log *slog.Logger
}

func newRepoPgx(l *slog.Logger) *repoPgx {
	name := slog.String("name", "poolRepoPgx")
	return &repoPgx{l.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return &repoPgx{}
}

func (r *repoPgx) Insert(source data.Source, root Root) (err error) {
	ds := data.MustConform[data.SourcePgx](source)
	dto := DataFromRoot(root)
	args := pgx.NamedArgs{
		"pool_id":     dto.PoolID,
		"title":       dto.Title,
		"sup_pool_id": dto.SupID,
		"revs":        dto.Revs,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRoot, args)
	if err != nil {
		r.log.Error("execution failed", slog.String("q", insertRoot))
		return err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "entity insertion succeeded", slog.Any("dto", dto))
	return nil
}

func (r *repoPgx) SelectAssets(source data.Source, poolID id.ADT) (AssetSnap, error) {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", poolID)
	rows, err := ds.Conn.Query(ds.Ctx, selectAssetSnap, poolID.String())
	if err != nil {
		r.log.Error("execution failed", idAttr, slog.String("q", selectAssetSnap))
		return AssetSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[assetSnapData])
	if err != nil {
		r.log.Error("row collection failed", idAttr)
		return AssetSnap{}, err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "entity selection succeeded", slog.Any("dto", dto))
	return DataToAssetSnap(dto)
}

func (r *repoPgx) SelectProc(source data.Source, procID id.ADT) (proc.Snap, error) {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("procID", procID)
	epRows, err := ds.Conn.Query(ds.Ctx, selectEPs, procID.String())
	if err != nil {
		r.log.Error("execution failed", idAttr, slog.String("q", selectEPs))
		return proc.Snap{}, err
	}
	defer epRows.Close()
	epDtos, err := pgx.CollectRows(epRows, pgx.RowToStructByName[epData])
	if err != nil {
		r.log.Error("rows collection failed", idAttr)
		return proc.Snap{}, err
	}
	eps, err := DataToEPs(epDtos)
	if err != nil {
		r.log.Error("mapping failed", idAttr)
		return proc.Snap{}, err
	}
	stepRows, err := ds.Conn.Query(ds.Ctx, selectSteps, procID.String())
	if err != nil {
		r.log.Error("execution failed", idAttr, slog.String("q", selectSteps))
		return proc.Snap{}, err
	}
	defer stepRows.Close()
	stepDtos, err := pgx.CollectRows(stepRows, pgx.RowToStructByName[step.RootData])
	if err != nil {
		r.log.Error("rows collection failed", idAttr)
		return proc.Snap{}, err
	}
	steps, err := step.DataToRoots(stepDtos)
	if err != nil {
		r.log.Error("mapping failed", idAttr)
		return proc.Snap{}, err
	}
	r.log.Debug("selection succeeded", idAttr)
	return proc.Snap{
		ProcID: procID,
		EPs:    core.IndexBy(proc.ProcPH, eps),
		Steps:  core.IndexBy(step.ChnlID, steps),
	}, nil
}

func (r *repoPgx) UpdateProc(source data.Source, mod proc.Mod) (err error) {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", mod.PoolID)
	dto := proc.DataFromMod(mod)
	// bindings
	bndReq := pgx.Batch{}
	for _, bnd := range dto.Bnds {
		args := pgx.NamedArgs{
			"proc_id":  bnd.ProcID,
			"proc_ph":  bnd.ProcPH,
			"chnl_id":  bnd.ChnlID,
			"state_id": bnd.StateID,
			"rev":      bnd.Rev,
		}
		bndReq.Queue(insertBnd, args)
	}
	bndRes := ds.Conn.SendBatch(ds.Ctx, &bndReq)
	defer func() {
		err = errors.Join(err, bndRes.Close())
	}()
	for _, bnd := range dto.Bnds {
		_, err := bndRes.Exec()
		if err != nil {
			r.log.Error("execution failed", idAttr, slog.String("q", insertBnd), slog.Any("bnd", bnd))
		}
	}
	if err != nil {
		return err
	}
	// steps
	stepReq := pgx.Batch{}
	for _, st := range dto.Steps {
		args := pgx.NamedArgs{
			"proc_id": st.PID,
			"chnl_id": st.VID,
			"kind":    st.K,
			"spec":    st.Spec,
		}
		stepReq.Queue(insertStep, args)
	}
	stepRes := ds.Conn.SendBatch(ds.Ctx, &stepReq)
	defer func() {
		err = errors.Join(err, bndRes.Close())
	}()
	for _, st := range dto.Steps {
		_, err := stepRes.Exec()
		if err != nil {
			r.log.Error("execution failed", idAttr, slog.String("q", insertStep), slog.Any("step", st))
		}
	}
	if err != nil {
		return err
	}
	// root
	args := pgx.NamedArgs{
		"pool_id": dto.PoolID,
		"rev":     dto.Rev,
		"k":       procRev,
	}
	ct, err := ds.Conn.Exec(ds.Ctx, updateRoot, args)
	if err != nil {
		r.log.Error("execution failed", idAttr, slog.String("q", updateRoot))
		return err
	}
	if ct.RowsAffected() == 0 {
		r.log.Error("update failed", idAttr)
		return errOptimisticUpdate(mod.Rev)
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "update succeeded", idAttr)
	return nil
}

func (r *repoPgx) UpdateAssets(source data.Source, mod AssetMod) (err error) {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", mod.OutPoolID)
	dto := DataFromAssetMod(mod)
	batch := pgx.Batch{}
	for _, ep := range dto.EPs {
		args := pgx.NamedArgs{
			"pool_id": dto.InPoolID,
			// proc_id ???
			"chnl_id":    ep.ChnlID,
			"state_id":   ep.StateID,
			"ex_pool_id": dto.OutPoolID,
			"rev":        dto.Rev,
		}
		batch.Queue(insertAsset, args)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, ep := range dto.EPs {
		_, err := br.Exec()
		if err != nil {
			r.log.Error("execution failed", idAttr, slog.String("q", insertAsset), slog.Any("ep", ep))
		}
	}
	if err != nil {
		return err
	}
	args := pgx.NamedArgs{
		"pool_id": dto.OutPoolID,
		"rev":     dto.Rev,
	}
	ct, err := ds.Conn.Exec(ds.Ctx, updateRoot, args)
	if err != nil {
		r.log.Error("execution failed", idAttr, slog.String("q", updateRoot))
		return err
	}
	if ct.RowsAffected() == 0 {
		r.log.Error("update failed", idAttr)
		return errOptimisticUpdate(mod.Rev)
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "update succeeded", idAttr)
	return nil
}

func (r *repoPgx) SelectSubs(source data.Source, poolID id.ADT) (SubSnap, error) {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", poolID)
	rows, err := ds.Conn.Query(ds.Ctx, selectOrgSnap, poolID.String())
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", selectOrgSnap))
		return SubSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[subSnapData])
	if err != nil {
		r.log.Error("row collection failed", idAttr)
		return SubSnap{}, err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "entity selection succeeded", slog.Any("dto", dto))
	return DataToSubSnap(dto)
}

func (r *repoPgx) SelectEPsByProcID(source data.Source, procID id.ADT) ([]proc.EP, error) {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", procID)
	rows, err := ds.Conn.Query(ds.Ctx, selectEPs, procID.String())
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", selectEPs))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[epData])
	if err != nil {
		r.log.Error("rows collection failed", idAttr)
		return nil, err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "entity selection succeeded", slog.Any("dtos", dtos))
	return DataToEPs(dtos)
}

func (r *repoPgx) SelectRefs(source data.Source) ([]Ref, error) {
	ds := data.MustConform[data.SourcePgx](source)
	query := `
		select
			pool_id, title
		from pool_roots`
	rows, err := ds.Conn.Query(ds.Ctx, query)
	if err != nil {
		r.log.Error("query execution failed", slog.String("q", query))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[refData])
	if err != nil {
		r.log.Error("rows collection failed")
		return nil, err
	}
	return DataToRefs(dtos)
}

func (r *repoPgx) Transfer(source data.Source, giver id.ADT, taker id.ADT, pids []chnl.ID) (err error) {
	ds := data.MustConform[data.SourcePgx](source)
	query := `
		insert into consumers (
			giver_id, taker_id, chnl_id
		) values (
			@giver_id, @taker_id, @chnl_id
		)`
	batch := pgx.Batch{}
	for _, pid := range pids {
		args := pgx.NamedArgs{
			"giver_id": sql.NullString{String: giver.String(), Valid: !giver.IsEmpty()},
			"taker_id": taker.String(),
			"chnl_id":  pid.String(),
		}
		batch.Queue(query, args)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, pid := range pids {
		_, err := br.Exec()
		if err != nil {
			r.log.Error("query execution failed",
				slog.Any("reason", err),
				slog.Any("id", pid),
			)
		}
	}
	if err != nil {
		return err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "context transfer succeeded")
	return nil
}

const (
	insertRoot = `
		insert into pool_roots (
			pool_id, revs, title, sup_pool_id
		) values (
			@pool_id, @revs, @title, @sup_pool_id
		)`

	insertAsset = `
		insert into pool_assets (
			pool_id, chnl_key, state_id, ex_pool_id, rev
		) values (
			@pool_id, @chnl_key, @state_id, @ex_pool_id, @rev
		)`

	insertBnd = `
		insert into pool_assets (
			pool_id, chnl_key, state_id, ex_pool_id, rev
		) values (
			@pool_id, @chnl_key, @state_id, @ex_pool_id, @rev
		)`

	insertStep = `
		insert into pool_steps (
			proc_id, chnl_id, kind, spec
		) values (
			@proc_id, @chnl_id, @kind, @spec
		)`

	updateRoot = `
		update pool_roots
		set revs[@k] = @rev + 1
		where pool_id = @pool_id
			and revs[@k] = @rev`

	selectOrgSnap = `
		select
			sup.pool_id,
			sup.title,
			jsonb_agg(json_build_object('pool_id', sub.pool_id, 'title', sub.title)) as subs
		from pool_roots sup
		left join pool_sups rel
			on rel.sup_pool_id = sup.pool_id
		left join pool_roots sub
			on sub.pool_id = rel.pool_id
			and sub.revs[1] = rel.rev
		where sup.pool_id = $1
		group by sup.pool_id, sup.title`

	selectAssetSnap = `
		select
			r.pool_id,
			r.title,
			jsonb_agg(json_build_object('chnl_key', a.chnl_key, 'state_id', a.state_id)) as ctx
		from pool_roots r
		left join pool_assets a
			on a.pool_id = r.pool_id
			and a.rev = r.revs[2]
		where r.pool_id = $1
		group by r.pool_id, r.title`

	selectEPs = `
		with bnds as not materialized (
			select distinct on (proc_ph)
				*
			from proc_bnds
			where proc_id = 'proc1'
			order by proc_ph, rev desc
		), liabs as not materialized (
			select distinct on (proc_ph)
				*
			from pool_liabs
			where proc_id = 'proc1'
			order by proc_ph, rev desc
		), assets as not materialized (
			select distinct on (proc_ph)
				*
			from pool_assets
			where proc_id = 'proc1'
			order by proc_ph, rev desc
		)
		select
			bnd.*,
			srv.pool_id as srv_id,
			srv.revs as srv_revs,
			clnt.pool_id as clnt_id,
			clnt.revs as clnt_revs
		from bnds bnd
		left join liabs liab
			on liab.proc_id = bnd.proc_id
			and liab.proc_ph = bnd.proc_ph
		left join pool_roots srv
			on srv.pool_id = liab.pool_id
		left join assets asset
			on asset.proc_id = bnd.proc_id
			and asset.proc_ph = bnd.proc_ph
		left join pool_roots clnt
			on clnt.pool_id = asset.pool_id`

	selectSteps = ``
)
