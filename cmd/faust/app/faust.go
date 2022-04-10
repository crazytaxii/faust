package faust

import (
	"context"

	"github.com/crazytaxii/faust/cmd/faust/app/config"
	"github.com/crazytaxii/faust/cmd/faust/app/options"
	"github.com/crazytaxii/faust/pkg/uploader"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewFaustApp(ver string) *cli.App {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	opts := options.NewAppOptions()
	cfg := config.NewAppConfig()

	return &cli.App{
		Name:    "faust",
		Usage:   "A simple tool for uploading image to object storage service",
		Version: ver,
		Action: func(c *cli.Context) error {
			return nil
		},
		Flags: opts.Flags(),
		Commands: []*cli.Command{
			{
				Name:    "upload",
				Aliases: []string{"up"},
				Usage:   "upload image to object storage service",
				Action: func(c *cli.Context) error {
					var err error
					defer func() {
						if err != nil {
							log.Errorf("error uploading image: %v", err)
						}
					}()

					cfg, err = opts.Config()
					if err != nil {
						return err
					}

					up := uploader.NewKodoUploader(cfg.QServiceConfig)
					ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
					defer cancel()
					res, err := up.Upload(ctx, opts.ImagePath)
					if err != nil {
						return err
					}
					log.Infof("image url: %v", res.URLs)
					return nil
				},
				Flags: cfg.QServiceConfig.Flags(),
			},
		},
	}
}
