package core

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("core",
	fx.Provide(
		fx.Annotate(newKeeper, fx.As(new(Keeper))),
	),
)

func newKeeper() *keeperViper {
	viper := viper.New()
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("reference")
	viper.ReadInConfig()
	viper.SetConfigName("application")
	viper.MergeInConfig()
	return &keeperViper{viper}
}
