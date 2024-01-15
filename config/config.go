package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App       `yaml:"app"`
		HTTP      `yaml:"http"`
		Log       `yaml:"logger"`
		Firestore `yaml:"firestore"`
		Redis     `yaml:"redis"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Firestore struct {
		ServiceAccLocation string `env-required:"true" yaml:"service_acc_key" `
		ProjectId          string `env-required:"true" yaml:"projectId"`
	}

	Redis struct {
		Address  string `env-required:"true" yaml:"server_address" `
		Password string `env-required:"true" yaml:"password" env:"REDIS_PASSWORD"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	// parentTop := filepath.Dir(wd)
	// parentTopTop := filepath.Dir(parentTop)
	// err = cleanenv.ReadConfig("./config/config.yml", cfg)
	// err = cleanenv.ReadConfig(parentTopTop+"/config/config.yml", cfg)
	err = cleanenv.ReadConfig(wd+"/config/config.yml", cfg)

	if err != nil {
		return nil, err
	}

	// err = cleanenv.ReadConfig(parentTopTop+"/.env", cfg)

	err = cleanenv.ReadConfig(wd+"/.env", cfg)

	// err = cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
