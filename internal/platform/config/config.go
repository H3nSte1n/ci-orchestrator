package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	App              AppConfig        `mapstructure:"app"`
	ApiServiceConfig ApiServiceConfig `mapstructure:"api_service"`
	DB               DBConfig         `mapstructure:"db"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ApiServiceConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func LoadConfig(env string) (config *Config, err error) {
	v := viper.New()
	v.SetConfigName(fmt.Sprintf(".env.%s", env))

	configDir := os.Getenv("BASE_DIR")
	if configDir == "" {
		configDir = "."
	}
	v.AddConfigPath(configDir)

	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return config, nil
}
