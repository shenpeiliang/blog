package config

import (
	"github.com/spf13/viper"
)

type (
	BaseConfig struct {
		Host      string `mapstructure:"host"`
		APIKey    string `mapstructure:"api_key"`
		MasterKey string `mapstructure:"master_key"`
	}

	IndexConfig struct {
		SearchableAttributes []string `mapstructure:"searchable_attributes"`
		FilterableAttributes []string `mapstructure:"filterable_attributes"`
		SortableAttributes   []string `mapstructure:"sortable_attributes"`
	}
	MeiliSearch struct {
		BaseConfig `mapstructure:",squash"`
		Indices    map[string]IndexConfig `mapstructure:"indices"`
	}

	Config struct {
		MeiliSearch `mapstructure:"meilisearch"`
	}
)

func LoadConfig() (*Config, error) {
	viper.SetConfigName("meilisearch")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
