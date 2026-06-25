package faust

import "github.com/urfave/cli/v3"

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
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Aliases:     []string{"c"},
			Usage:       "`file` path of configuration to be loaded",
			Destination: &o.ConfigFile,
		},
		&cli.StringFlag{
			Name:        "image",
			Aliases:     []string{"i"},
			Usage:       "`file` path of image to be uploaded",
			Destination: &o.ImagePath,
		},
		&cli.StringFlag{
			Name:        "cert",
			Usage:       "`file` path of certificate to be uploaded",
			Destination: &o.CertPath,
		},
		&cli.StringFlag{
			Name:        "key",
			Usage:       "`file` path of private key to be uploaded",
			Destination: &o.KeyPath,
		},
	}
}

func (o *Options) LoadConfig() error {
	// load config file first
	cfg, err := TryToLoadConfig(o.ConfigFile)
	if err != nil {
		return err
	}

	o.Config = cfg
	return nil
}
