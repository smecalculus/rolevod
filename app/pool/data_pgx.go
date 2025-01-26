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
	insertRoot := `
		insert into pool_roots (
			pool_id, rev, title, sup_id
		) values (
			@pool_id, @rev, @title, @sup_id
		)`
	rootArgs := pgx.NamedArgs{
		"pool_id": dto.ID,
		"rev":     dto.Rev,
		"title":   dto.Title,
		"sup_id":  dto.SupID,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRoot, rootArgs)
	if err != nil {
		r.log.Error("query execution failed", slog.String("q", insertRoot))
		return err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "entity insertion succeeded", slog.Any("dto", dto))
	return nil
}

func (r *repoPgx) SelectByID(source data.Source, rid id.ADT) (Snap, error) {
	ds := data.MustConform[data.SourcePgx](source)
	idAttr := slog.Any("id", rid)
	rows, err := ds.Conn.Query(ds.Ctx, selectById, rid.String())
	if err != nil {
		r.log.Error("query execution failed", idAttr, slog.String("q", selectById))
		return Snap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[snapData])
	if err != nil {
		r.log.Error("row collection failed", idAttr, slog.Any("reason", err))
		return Snap{}, err
	}
	r.log.Log(ds.Ctx, core.LevelTrace, "entity selection succeeded", slog.Any("dto", dto))
	return DataToSnap(dto)
}

func (r *repoPgx) SelectAll(source data.Source) ([]Ref, error) {
	ds := data.MustConform[data.SourcePgx](source)
	query := `
		select
			pool_id, rev, title
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
	selectById = `
		select
			tr.pool_id,
			(array_agg(tr.title))[1] as title,
			jsonb_agg(to_jsonb((select sub from (select ts.pool_id, ts.rev, ts.title) sub))) filter (where ts.pool_id is not null) as subs
		from pool_roots tr
		left join pool_roots ts
			on ts.sup_id = tr.pool_id
		where tr.pool_id = $1
		group by tr.pool_id`
)
