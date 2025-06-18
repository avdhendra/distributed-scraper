package config

import (
	"bytes"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

type OAuthConfig struct {
	LinkedInClientID     string `mapstructure:"LINKEDIN_CLIENT_ID"`
	LinkedInClientSecret string `mapstructure:"LINKEDIN_CLIENT_SECRET"`
	YouTubeClientID      string `mapstructure:"YOUTUBE_CLIENT_ID"`
	YouTubeClientSecret  string `mapstructure:"YOUTUBE_CLIENT_SECRET"`
	InstagramClientID    string `mapstructure:"INSTAGRAM_CLIENT_ID"`
	InstagramClientSecret string `mapstructure:"INSTAGRAM_CLIENT_SECRET"`
}

type Config struct {
	KafkaBrokers   []string    `mapstructure:"KAFKA_BROKERS"`
	ScrapeInterval int         `mapstructure:"SCRAPE_INTERVAL"`
	ProxyList      []string    `mapstructure:"PROXY_LIST"`
	OAuthConfig    OAuthConfig `mapstructure:"OAUTH"`
}

func LoadFromConsul() (*Config, error) {
	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.SetConfigFile(".env")
	v.AutomaticEnv()
	v.SetDefault("SCRAPE_INTERVAL", 300)

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