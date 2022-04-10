package uploader

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/crazytaxii/faust/pkg/uploader/utils"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/urfave/cli/v2"
	_ "golang.org/x/image/webp"
)

const (
	defaultExpires    = 3600
	fmtKodoReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`
)

type Uploader interface {
	Upload() error
}

type (
	// qiniu
	QServiceConfig struct {
		AccessKey string `jaon:"access_key" yaml:"accessKey"`
		SecretKey string `json:"secret_key" yaml:"secretKey"`
		Expires   uint64 `json:"expires" yaml:"expires"`
		Bucket    string `json:"bucket" yaml:"bucket"`
	}
	// KodoUploader
	KodoUploader struct {
		config        *QServiceConfig
		credentials   *qbox.Mac
		uploader      *storage.FormUploader
		bucketManager *storage.BucketManager
	}
	KodoPutRet struct {
		Key    string `json:"key"`
		Hash   string `json:"hash"`
		Fsize  uint64 `json:"fsize"`
		Bucket string `json:"bucket"`
		URLs   []string
	}
)

func NewQServiceConfig() *QServiceConfig {
	return &QServiceConfig{
		Expires: defaultExpires,
	}
}

func (c *QServiceConfig) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "access-key",
			Aliases:     []string{"a"},
			Usage:       "access key of qiniu kodo object storage service",
			Destination: &c.AccessKey,
		},
		&cli.StringFlag{
			Name:        "secret-key",
			Aliases:     []string{"s"},
			Usage:       "secret key of qiniu kodo object storage service",
			Destination: &c.SecretKey,
		},
		&cli.Uint64Flag{
			Name:        "expires",
			Aliases:     []string{"e"},
			Usage:       "expires time",
			Destination: &c.Expires,
		},
		&cli.StringFlag{
			Name:        "bucket",
			Aliases:     []string{"b"},
			Usage:       "bucket",
			Destination: &c.Bucket,
		},
	}
}

func (c *QServiceConfig) MakePutPolicy() *storage.PutPolicy {
	return &storage.PutPolicy{
		Scope:      c.Bucket,
		Expires:    c.Expires,
		ReturnBody: fmtKodoReturnBody,
	}
}

func NewKodoUploader(cfg *QServiceConfig) *KodoUploader {
	sc := &storage.Config{}
	mac := auth.New(cfg.AccessKey, cfg.SecretKey)
	return &KodoUploader{
		config:        cfg,
		credentials:   mac,
		uploader:      storage.NewFormUploader(sc),
		bucketManager: storage.NewBucketManager(mac, sc),
	}
}

// Upload reads file and put on to specific bucket of Kodo.
func (u *KodoUploader) Upload(ctx context.Context, name string) (*KodoPutRet, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	_, format, err := image.DecodeConfig(f)
	if err != nil {
		return nil, err
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	pp := u.config.MakePutPolicy()
	token := pp.UploadToken(u.credentials)
	ret := &KodoPutRet{}
	key := utils.GenUploadKey(format)
	// upload file
	if err := u.uploader.Put(ctx, ret, token, key, f, fi.Size(), &storage.PutExtra{}); err != nil {
		return nil, err
	}

	// query bucket domains
	domainInfo, err := u.bucketManager.ListBucketDomains(u.config.Bucket)
	if err != nil {
		return nil, err
	}
	for _, domain := range domainInfo {
		ret.URLs = append(ret.URLs, fmt.Sprintf("https://%s/%s", domain.Domain, ret.Key))
	}

	return ret, nil
}
