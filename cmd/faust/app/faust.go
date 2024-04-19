package faust

import (
	"context"

	"github.com/crazytaxii/faust/cmd/faust/app/config"
	"github.com/crazytaxii/faust/cmd/faust/app/options"
	"github.com/crazytaxii/faust/pkg/uploader"

	"github.com/docker/go-units"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func NewFaustApp(ver string) *cli.App {
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
				Action: func(c *cli.Context) (err error) {
					defer func() {
						if err != nil {
							log.Errorf("error uploading image: %v", err)
						}
					}()

					if cfg, err = opts.Config(); err != nil {
						return err
					}

					up := uploader.NewKodoUploader(cfg.QServiceConfig)
					ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
					defer cancel()
					res, err := up.Upload(ctx, opts.ImagePath)
					if err != nil {
						return err
					}
					log.WithFields(
						log.Fields{
							"key":  res.Key,
							"size": units.HumanSize(float64(res.Fsize)),
							"hash": res.Hash,
						},
					).Infof("image url: %v", res.URLs)
					return nil
				},
				Flags: cfg.QServiceConfig.Flags(),
			},
		},
	}
}
