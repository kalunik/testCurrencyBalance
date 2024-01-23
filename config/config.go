package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

type AppConfig struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Nats     NatsConfig
}

type ServerConfig struct {
	Port              string
	Ssl               bool
	CtxDefaultTimeout time.Duration
}

type PostgresConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DbName       string
	SSLMode      string
	PoolMaxConns int
}

type NatsConfig struct {
	Host              string
	Port              string
	Cluster           string
	Client            string
	Subject           string
	DurableName       string
	ReconnectInterval time.Duration
	ServerPort        string
}

type Config interface {
	ParseConfig() (AppConfig, error)
}

type ConfigDriver struct {
	v *viper.Viper
}

func LoadNewConfig() (Config, error) {
	v := viper.New()
	v.SetConfigFile(findConfigPath())

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("config: file not found")
		}
		return nil, err
	}

	return &ConfigDriver{
		v: v,
	}, nil
}

func (c *ConfigDriver) ParseConfig() (AppConfig, error) {
	config := AppConfig{}

	err := c.v.Unmarshal(&config)
	if err != nil {
		return AppConfig{}, fmt.Errorf("config: unable to decode into struct: %w", err)
	}
	return config, nil
}

func findConfigPath() string {
	customPath := os.Getenv("CONFIG_PATH")

	configPaths := map[string]string{
		"docker": "config/config_docker.yml",
		"local":  "config/config_local.yml",
		"custom": customPath,
	}
	defaultPath := configPaths["local"]

	pathKey := os.Getenv("CONFIG")
	if pathKey == "" {
		return defaultPath
	}
	if path, ok := configPaths[pathKey]; ok {
		return path
	}
	return defaultPath
}
