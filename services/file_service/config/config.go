package config

import (
	"github.com/spf13/viper"
)

type Model struct{}

type AppConfigInterface interface {
	InitializeConfig(string) error
}

func AppConfig() *Model {
	return &Model{}
}

type ProjectConfig struct {
	BundleId          string   `mapStructure:"bundleId"`
	AcceptedFileTypes []string `mapStructure:"acceptedFileTypes"`
}

var allConfig ProjectConfig

func (model *Model) InitializeConfig(configPath string) error {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&allConfig)
}

func (model *Model) GetAllConfig() *ProjectConfig {
	return &allConfig
}
