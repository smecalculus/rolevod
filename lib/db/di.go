package db

import (
	"context"
	"smecalculus/rolevod/lib/cfg"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module("db",
	fx.Provide(
		newCfg,
		newPgx,
	),
)

func newCfg(keeper cfg.Keeper) (*props, error) {
	props := &props{}
	err := keeper.Load("db", props)
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
