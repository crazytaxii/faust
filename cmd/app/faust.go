package faust

import (
	"context"
	"fmt"

	"github.com/crazytaxii/faust/pkg/service"

	"github.com/docker/go-units"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func NewFaustApp(ver string) *cli.Command {
	gopts := NewGlobalOptions()
	return &cli.Command{
		Name:    "faust",
		Usage:   "A simple tool for uploading image to object storage service",
		Version: ver,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return cli.ShowAppHelp(cmd)
		},
		Commands: []*cli.Command{
			uploadCmd(gopts),
			deleteCmd(gopts),
		},
		Flags: gopts.Flags(),
	}
}

func uploadCmd(gopts *GlobalOptions) *cli.Command {
	opts := gopts.NewUploadOptions()
	return &cli.Command{
		Name:    "upload",
		Aliases: []string{"up"},
		Usage:   "Upload image or certificates to object storage service",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := runUpload(ctx, cmd, opts); err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
		Flags: opts.Flags(),
	}
}

func runUpload(ctx context.Context, cmd *cli.Command, opts *UploadOptions) error {
	cfg, err := opts.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	si := service.NewQiniuService(cfg.QServiceConfig)
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
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
		return cli.ShowSubcommandHelp(cmd)
	}

	log.WithFields(lf).Info("upload successfully")
	return nil
}

func deleteCmd(gopts *GlobalOptions) *cli.Command {
	opts := gopts.NewDeleteOptions()
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"del"},
		Usage:   "Delete image from object storage service",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := runDelete(ctx, cmd, opts); err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
		Flags: opts.Flags(),
	}
}

func runDelete(ctx context.Context, cmd *cli.Command, opts *DeleteOptions) error {
	cfg, err := opts.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	si := service.NewQiniuService(cfg.QServiceConfig)
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	if opts.Key != "" {
		if err := si.DeleteImage(ctx, opts.Key); err != nil {
			return fmt.Errorf("error deleting image: %w", err)
		}
	} else {
		return cli.ShowSubcommandHelp(cmd)
	}

	log.WithFields(log.Fields{
		"key":    opts.Key,
		"bucket": cfg.QServiceConfig.Bucket,
	}).Info("delete successfully")
	return nil
}
