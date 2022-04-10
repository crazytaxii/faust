package options

import (
	"github.com/crazytaxii/faust/cmd/faust/app/config"
	"github.com/urfave/cli/v2"
)

type AppOptions struct {
	// path of config file
	ConfigFile string
	// path of image to be upload
	ImagePath string
}

func NewAppOptions() *AppOptions {
	return &AppOptions{}
}

func (o *AppOptions) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Aliases:     []string{"c"},
			Usage:       "file path of configuration to be loaded",
			Destination: &o.ConfigFile,
		},
		&cli.StringFlag{
			Name:        "image",
			Aliases:     []string{"i"},
			Usage:       "file path of image to be uploaded",
			Destination: &o.ImagePath,
		},
	}
}

func (o *AppOptions) Config() (*config.AppConfig, error) {
	// try to load config file
	return config.TryToLoadConfig(o.ConfigFile)
}
