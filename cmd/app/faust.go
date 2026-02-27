package faust

import (
	"context"
	"fmt"

	"github.com/crazytaxii/faust/pkg/service"

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
	opts := NewOptions()

	return &cli.App{
		Name:    "faust",
		Usage:   "A simple tool for uploading image to object storage service",
		Version: ver,
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			{
				Name:    "upload",
				Aliases: []string{"up"},
				Usage:   "Upload image or certificates to object storage service",
				Action: func(c *cli.Context) error {
					if err := runUpload(c, opts); err != nil {
						log.Error(err)
						return err
					}
					return nil
				},
				Flags: opts.Flags(),
			},
		},
	}
}

func runUpload(c *cli.Context, opts *Options) error {
	if err := opts.LoadConfig(c); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	var si service.ServiceInterface = service.NewQiniuService(opts.Config.QServiceConfig)
	ctx, cancel := context.WithTimeout(c.Context, opts.Config.Timeout)
	defer cancel()

	lf := make(log.Fields)
	if opts.ImagePath != "" {
		res, err := si.UploadImage(ctx, opts.ImagePath)
		if err != nil {
			return fmt.Errorf("error uploading image: %w", err)
		}
		lf["bucket"] = res.Bucket
		lf["key"] = res.Key
		lf["size"] = units.HumanSize(float64(res.Size))
		lf["image_url"] = res.URLs
	} else if opts.CertPath != "" && opts.KeyPath != "" {
		res, err := si.UploadCerts(ctx, opts.KeyPath, opts.CertPath)
		if err != nil {
			return fmt.Errorf("error uploading certificates: %w", err)
		}
		lf["common_name"] = res.CommonName
		lf["expiration"] = res.Expiration
	} else {
		return cli.ShowSubcommandHelp(c)
	}

	log.WithFields(lf).Info("upload successfully")
	return nil
}
