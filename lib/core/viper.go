package core

import (
	"github.com/spf13/viper"
)

// adapter
type keeperViper struct {
	viper *viper.Viper
}

func (k *keeperViper) Load(key string, v any) error {
	return k.viper.UnmarshalKey(key, v)
}
