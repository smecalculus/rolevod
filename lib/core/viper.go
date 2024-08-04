package core

import (
	"log/slog"

	"github.com/spf13/viper"
)

// adapter
type keeperViper struct {
	viper  *viper.Viper
	logger *slog.Logger
}

func (k *keeperViper) Load(key string, v any) error {
	err := k.viper.UnmarshalKey(key, v)
	if err != nil {
		return err
	}
	k.logger.Info("loaded", slog.Any("props", v))
	return nil
}
