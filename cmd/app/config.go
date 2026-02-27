package faust

import (
	"errors"
	"os"
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

type Config struct {
	*service.QServiceConfig `mapstructure:",squash"`
	Timeout                 time.Duration `json:"timeout" yaml:"timeout"`
}

func NewConfig() *Config {
	return &Config{
		QServiceConfig: service.NewQServiceConfig(),
		Timeout:        defaultTimeout,
	}
}

func TryToLoadConfig(file string) (*Config, error) {
	cfg := NewConfig()
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

	v := viper.New()
	v.SetConfigType(strings.TrimPrefix(ext, "."))
	v.SetConfigName(strings.TrimSuffix(fileName, ext))
	v.AddConfigPath(cfgPath)
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	v.AddConfigPath(filepath.Join(home, ".faust"))
	v.AddConfigPath("/etc/faust")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return cfg, nil
		}
		return nil, err
	}

	if err := v.Unmarshal(cfg); err != nil {
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
