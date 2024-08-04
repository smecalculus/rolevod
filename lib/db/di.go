package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
)

var Module = fx.Module("db",
	fx.Provide(
		newPgx,
	),
	fx.Provide(
		fx.Private,
		newCfg,
	),
)

func newCfg(k core.Keeper) (*props, error) {
	props := &props{}
	err := k.Load("db", props)
	if err != nil {
		return nil, err
	}
	return props, nil
}

func newPgx(props *props, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	pgx, err := pgxpool.New(context.Background(), props.Protocol.Postgres.Url)
	if err != nil {
		return nil, err
	}
	lc.Append(
		fx.Hook{
			OnStart: pgx.Ping,
			OnStop: func(ctx context.Context) error {
				go pgx.Close()
				return nil
			},
		},
	)
	return pgx, nil
}
