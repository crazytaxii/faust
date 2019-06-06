package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	APP_NAME    = "faust"
	APP_USAGE   = "A simple tool for qiniu cloud"
	APP_VERSION = "0.9.1"
)

type Config struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
	BaseUrl   string `json:"base_url"`
}

var conf *Config
var imgPath string

const (
	configFile = "./config.json"
)

func init() {
	conf = &Config{}
	err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = APP_NAME
	app.Usage = APP_USAGE
	app.Version = APP_VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "image",
			Usage:       "image path",
			Destination: &imgPath,
		},
		cli.StringFlag{
			Name:        "access_key",
			Usage:       "access key",
			Value:       conf.AccessKey,
			Destination: &conf.AccessKey,
		},
		cli.StringFlag{
			Name:        "secret_key",
			Usage:       "secret key",
			Value:       conf.SecretKey,
			Destination: &conf.SecretKey,
		},
		cli.StringFlag{
			Name:        "bucket",
			Usage:       "bucket name",
			Value:       conf.Bucket,
			Destination: &conf.Bucket,
		},
		cli.StringFlag{
			Name:        "domain",
			Usage:       "your domain",
			Value:       conf.Domain,
			Destination: &conf.Domain,
		},
		cli.StringFlag{
			Name:        "base_url",
			Usage:       "base url",
			Value:       conf.BaseUrl,
			Destination: &conf.BaseUrl,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "upload",
			Aliases: []string{"u"},
			Usage:   "Uploads image to Object Storage",
			Action:  upload,
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Saves configurations",
			Action:  saveConfig,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
