package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		DatabaseURL string `mapstructure:"DATABASE_URL"`
		ScaleFactor int    `mapstructure:"SCALE_FACTOR"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	return &cfg, err
}
