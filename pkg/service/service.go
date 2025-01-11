package service

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/crazytaxii/faust/pkg/image"
	"github.com/crazytaxii/faust/pkg/service/utils"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/urfave/cli/v2"
	_ "golang.org/x/image/webp"
)

const (
	defaultExpires    = 3600
	fmtKodoReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`
)

type (
	ImageUploadResponse struct {
		Key  string
		Hash string
		Size uint64
		URLs []string
	}
	CertsUploadResponse struct {
		CommonName string
		Expiration time.Time
	}
)

type ServiceInterface interface {
	UploadImage(ctx context.Context, image string) (*ImageUploadResponse, error)
	UploadCerts(ctx context.Context, key, cert string) (*CertsUploadResponse, error)
}

type (
	QServiceConfig struct {
		AccessKey string `json:"access_key" yaml:"accessKey"`
		SecretKey string `json:"secret_key" yaml:"secretKey"`
		Expires   uint64 `json:"expires" yaml:"expires"`
		Bucket    string `json:"bucket" yaml:"bucket"`
	}
	QiniuService struct {
		config        *QServiceConfig
		credentials   *qbox.Mac
		uploader      *storage.FormUploader
		bucketManager *storage.BucketManager
	}
	kodoPutRet struct {
		Key    string `json:"key"`
		Hash   string `json:"hash"`
		Fsize  uint64 `json:"fsize"`
		Bucket string `json:"bucket"`
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

func NewQiniuService(cfg *QServiceConfig) *QiniuService {
	sc := &storage.Config{}
	mac := auth.New(cfg.AccessKey, cfg.SecretKey)
	return &QiniuService{
		config:        cfg,
		credentials:   mac,
		uploader:      storage.NewFormUploader(sc),
		bucketManager: storage.NewBucketManager(mac, sc),
	}
}

// UploadImage reads file and put on to the specific bucket.
func (s *QiniuService) UploadImage(ctx context.Context, name string) (*ImageUploadResponse, error) {
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

	pp := s.config.MakePutPolicy()
	token := pp.UploadToken(s.credentials)
	ret := &kodoPutRet{}
	key := utils.GenUploadKey(format)
	// upload file
	if err := s.uploader.Put(ctx, ret, token, key, f, fi.Size(), &storage.PutExtra{}); err != nil {
		return nil, err
	}

	// query bucket domains
	domainInfo, err := s.bucketManager.ListBucketDomains(s.config.Bucket)
	if err != nil {
		return nil, err
	}
	urls := make([]string, len(domainInfo))
	for i, domain := range domainInfo {
		urls[i] = fmt.Sprintf("https://%s/%s", domain.Domain, ret.Key)
	}

	return &ImageUploadResponse{
		Key:  ret.Key,
		Hash: ret.Hash,
		Size: ret.Fsize,
		URLs: urls,
	}, nil
}

type (
	SSLCerts struct {
		Name       string `json:"name"`
		CommonName string `json:"common_name"`
		Key        string `json:"pri"`
		CertChain  string `json:"ca"`
	}
)

func (s *QiniuService) UploadCerts(ctx context.Context, keyPath, certPath string) (*CertsUploadResponse, error) {
	rawKeyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(rawKeyData)
	if block == nil {
		return nil, errors.New("invalid private key: no PEM data found")
	}

	rawCertData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	if block, _ = pem.Decode(rawCertData); block == nil {
		return nil, errors.New("invalid certificate: no PEM data found")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}
	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		return nil, errors.New("certificate is expired or not yet valid")
	}

	// doc: https://developer.qiniu.com/fusion/8593/interface-related-certificate
	reqBody := &SSLCerts{
		Name:       cert.Subject.CommonName,
		CommonName: cert.Subject.CommonName,
		Key:        string(rawKeyData),
		CertChain:  string(rawCertData),
	}
	if _, err = postRequest(ctx, s.credentials, "/sslcert", reqBody); err != nil {
		return nil, err
	}
	return &CertsUploadResponse{
		CommonName: cert.Subject.CommonName,
		Expiration: cert.NotAfter,
	}, nil
}

func postRequest(ctx context.Context, mac *auth.Credentials, path string, body interface{}) ([]byte, error) {
	reqData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s%s", storage.DefaultAPIHost, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}

	token, err := mac.SignRequest(req)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("QBox %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ResponseError(resp)
	}
	return io.ReadAll(resp.Body)
}
