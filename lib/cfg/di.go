package cfg

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Module("cfg",
	fx.Provide(
		fx.Annotate(newKeeper, fx.As(new(Keeper))),
	),
)

func newKeeper() *keeperViper {
	viper := viper.New()
	viper.AddConfigPath(".")
	viper.SetConfigName("rolevod")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	return &keeperViper{viper}
}
