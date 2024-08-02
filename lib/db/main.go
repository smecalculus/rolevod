package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module("db",
	fx.Provide(
		newPgx,
	),
)

func newPgx(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	pgx, err := pgxpool.New(context.Background(), "")
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
