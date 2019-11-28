package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	AppName  = "faust"
	AppUsage = "A simple tool for qiniu cloud"
)

var AppVersion string

type QiniuConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	BaseUrl   string `yaml:"base_url"`
}

type FaustClient struct {
	Config         *QiniuConfig
	App            *cli.App
	ConfigFilePath string
}

const ConfigFile = "./config.yaml"

func main() {
	fc := &FaustClient{}
	err := fc.LoadConfig(ConfigFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fc.NewApp(AppName, AppUsage, AppVersion)
	fc.App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "image",
			Usage: "image path",
		},
		cli.StringFlag{
			Name:        "access_key",
			Usage:       "access key",
			Value:       fc.Config.AccessKey,
			Destination: &fc.Config.AccessKey,
		},
		cli.StringFlag{
			Name:        "secret_key",
			Usage:       "secret key",
			Value:       fc.Config.SecretKey,
			Destination: &fc.Config.SecretKey,
		},
		cli.StringFlag{
			Name:        "bucket",
			Usage:       "bucket name",
			Value:       fc.Config.Bucket,
			Destination: &fc.Config.Bucket,
		},
		cli.StringFlag{
			Name:        "base_url",
			Usage:       "base url",
			Value:       fc.Config.BaseUrl,
			Destination: &fc.Config.BaseUrl,
		},
	}
	fc.App.Commands = []cli.Command{
		{
			Name:    "upload",
			Aliases: []string{"u"},
			Usage:   "Uploads image to Object Storage",
			Action:  fc.Upload,
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Saves configurations",
			Action: func(c *cli.Context) error {
				return fc.SaveConfig()
			},
		},
	}
	err = fc.App.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// New a cli app
func (fc *FaustClient) NewApp(name, usage, version string) {
	fc.App = cli.NewApp()
	fc.App.Name = name
	fc.App.Usage = usage
	fc.App.Version = version
	fc.Config = &QiniuConfig{}
}
