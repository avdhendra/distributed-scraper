package config

import (
	"bytes"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresURL  string   `mapstructure:"POSTGRES_URL"`
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
}

func LoadFromConsul() (*Config, error) {
	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.SetConfigFile(".env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	client, err := api.NewClient(&api.Config{Address: "consul:8500"})
	if err != nil {
		return nil, err
	}

	kv := client.KV()
	pair, _, err := kv.Get("scraper/config", nil)
	if err != nil {
		return nil, err
	}

	if pair != nil {
		if err := v.ReadConfig(bytes.NewReader(pair.Value)); err != nil {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}