package config

import (
	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.AddConfigPath("billing/config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}
	viper.AutomaticEnv()
	return nil
}
