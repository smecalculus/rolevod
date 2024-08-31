package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"smecalculus/rolevod/lib/core"
)

var Module = fx.Module("lib/store",
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
	err := k.Load("storage", props)
	if err != nil {
		return nil, err
	}
	return props, nil
}

func newPgx(p *props, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(p.Protocol.Postgres.Url)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 2
	pgx, err := pgxpool.NewWithConfig(context.Background(), config)
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
