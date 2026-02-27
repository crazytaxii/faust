package faust

import (
	"github.com/crazytaxii/faust/pkg/service"

	"github.com/urfave/cli/v2"
)

type Options struct {
	// path of config file
	ConfigFile string
	// path of image to be upload
	Config *Config

	ImagePath string
	CertPath  string
	KeyPath   string
}

func NewOptions() *Options {
	return &Options{
		Config: NewConfig(),
	}
}

func (o *Options) Flags() []cli.Flag {
	return append([]cli.Flag{
		&cli.PathFlag{
			Name:        "config-file",
			Aliases:     []string{"c"},
			Usage:       "`file` path of configuration to be loaded",
			Destination: &o.ConfigFile,
		},
		&cli.PathFlag{
			Name:        "image",
			Aliases:     []string{"i"},
			Usage:       "`file` path of image to be uploaded",
			Destination: &o.ImagePath,
		},
		&cli.PathFlag{
			Name:        "cert",
			Usage:       "`file` path of certificate to be uploaded",
			Destination: &o.CertPath,
		},
		&cli.PathFlag{
			Name:        "key",
			Usage:       "`file` path of private key to be uploaded",
			Destination: &o.KeyPath,
		},
	}, o.Config.Flags()...)
}

func (o *Options) LoadConfig(c *cli.Context) error {
	// load config file first
	cfg, err := TryToLoadConfig(o.ConfigFile)
	if err != nil {
		return err
	}

	// override config file values with explicitly set CLI flags
	o.applyFlagOverrides(c, cfg.QServiceConfig)
	o.Config = cfg
	return nil
}

// applyFlagOverrides re-applies CLI flag values to the loaded config
func (o *Options) applyFlagOverrides(c *cli.Context, qCfg *service.QServiceConfig) {
	if c.IsSet("access-key") {
		qCfg.AccessKey = o.Config.AccessKey
	}
	if c.IsSet("secret-key") {
		qCfg.SecretKey = o.Config.SecretKey
	}
	if c.IsSet("expires") {
		qCfg.Expires = o.Config.Expires
	}
	if c.IsSet("bucket") {
		qCfg.Bucket = o.Config.Bucket
	}
}
