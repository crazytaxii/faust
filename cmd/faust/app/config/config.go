package config

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/crazytaxii/faust/pkg/service"

	"github.com/spf13/viper"
)

const (
	defaultConfigName = "config.yaml"
	defaultTimeout    = 10 * time.Second
)

type AppConfig struct {
	*service.QServiceConfig `mapstructure:",squash"`
	Timeout                 time.Duration `json:"timeout" yaml:"timeout"`
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		QServiceConfig: service.NewQServiceConfig(),
		Timeout:        defaultTimeout,
	}
}

func TryToLoadConfig(file string) (*AppConfig, error) {
	cfg := NewAppConfig()
	if file == "" {
		file = defaultConfigName
	}
	cfgFile, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	cfgPath := filepath.Dir(cfgFile)
	fileName := filepath.Base(cfgFile)
	ext := filepath.Ext(fileName)

	if err := validateExt(ext); err != nil {
		return nil, err
	}

	viper.SetConfigType(strings.TrimPrefix(ext, "."))
	viper.SetConfigName(strings.TrimSuffix(fileName, ext))
	viper.AddConfigPath(cfgPath)
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(path.Join(home, ".faust"))
	viper.AddConfigPath("/etc/faust")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return cfg, nil
		}
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validateExt(ext string) error {
	switch ext {
	case ".json", ".yaml", ".yml":
	default:
		return errors.New("unsupported config file extension")
	}
	return nil
}
