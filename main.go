package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/urfave/cli"
)

const (
	AppName  = "faust"
	AppUsage = "A simple tool for qiniu cloud"
)

var (
	AppVersion string
	ConfigFile string
)

type FaustClient struct {
	Config         *QiniuConfig
	App            *cli.App
	ConfigFilePath string
	ImgPath        string
}

func init() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	abs := fmt.Sprintf("%s/.faust", user.HomeDir)
	// check ~/.faust is existed
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		// mkdir ~/.faust
		if err := os.Mkdir(abs, 0755); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
	ConfigFile = fmt.Sprintf("%s/config.yaml", abs)
}

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
			Name:        "image",
			Usage:       "image path",
			Destination: &fc.ImgPath,
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

// New a cli app.
func (fc *FaustClient) NewApp(name, usage, version string) {
	fc.App = cli.NewApp()
	fc.App.Name = name
	fc.App.Usage = usage
	fc.App.Version = version
}
