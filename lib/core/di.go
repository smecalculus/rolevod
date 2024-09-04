package core

import (
	"log/slog"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("lib/core",
	fx.Provide(
		newLogger,
		fx.Annotate(newKeeper, fx.As(new(Keeper))),
	),
)

func newLogger() *slog.Logger {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	return slog.Default()
}

func newKeeper(l *slog.Logger) *keeperViper {
	viper := viper.New()
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("reference")
	viper.ReadInConfig()
	viper.SetConfigName("application")
	viper.MergeInConfig()
	t := slog.String("t", "core.keeperViper")
	return &keeperViper{viper, l.With(t)}
}
