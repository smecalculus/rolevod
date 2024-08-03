package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module("db",
	fx.Provide(
		newConf,
		newPgx,
	),
)

func newConf() props {
	return props{
		Protocol: protocol{
			Postgres: postgres{
				Url: "postgres://postgres:password@localhost:5432/postgres",
			},
		},
	}
}

func newPgx(props props, lc fx.Lifecycle) (*pgxpool.Pool, error) {
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
