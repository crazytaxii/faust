package faust

import (
	"context"

	"github.com/crazytaxii/faust/cmd/faust/app/config"
	"github.com/crazytaxii/faust/cmd/faust/app/options"
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
	opts := options.NewAppOptions()
	cfg := config.NewAppConfig()

	return &cli.App{
		Name:    "faust",
		Usage:   "A simple tool for uploading image to object storage service",
		Version: ver,
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Flags: opts.Flags(),
		Commands: []*cli.Command{
			{
				Name:    "upload",
				Aliases: []string{"up"},
				Usage:   "Upload image or certificates to object storage service",
				Action: func(c *cli.Context) error {
					return runUpload(c.Context, opts)
				},
				Flags: cfg.QServiceConfig.Flags(),
			},
		},
	}
}

func runUpload(ctx context.Context, opts *options.AppOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		log.Errorf("error loading config: %v", err)
		return err
	}
	var si service.ServiceInterface = service.NewQiniuService(cfg.QServiceConfig)
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	lf := make(log.Fields)
	if opts.ImagePath != "" {
		res, err := si.UploadImage(ctx, opts.ImagePath)
		if err != nil {
			log.Errorf("error uploading image: %v", err)
			return err
		}
		lf["key"] = res.Key
		lf["size"] = units.HumanSize(float64(res.Size))
		lf["hash"] = res.Hash
		lf["image_url"] = res.URLs
	}

	if opts.CertPath != "" && opts.KeyPath != "" {
		res, err := si.UploadCerts(ctx, opts.KeyPath, opts.CertPath)
		if err != nil {
			log.Errorf("error uploading certificates: %v", err)
		}
		if res != nil {
			log.Errorf("error uploading certificates: %v", err)
			return err
		}
		lf["common_name"] = res.CommonName
		lf["expiration"] = res.Expiration
	}

	log.WithFields(lf).Info("upload successfully")
	return nil
}
