package config

import "github.com/spf13/viper"

type Config struct {
	KafkaUri             string `mapstructure:"KAFKA_URI"`
	KafkaTopic           string `mapstructure:"KAFKA_TOPIC"`
	PsqlConnectionString string `mapstructure:"PSQL_CONNECTION_STRING"`
}

func NewConfig(path string, name string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
