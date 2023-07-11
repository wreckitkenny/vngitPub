package utils

import (
	"vngitPub/model"
	"github.com/spf13/viper"
)

// LoadConfig - load RabbitMQ configs
func LoadConfig(path string) (config model.Default, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("rabbitmq")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}