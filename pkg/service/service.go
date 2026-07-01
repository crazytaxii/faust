package service

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/crazytaxii/faust/pkg/image"
	"github.com/crazytaxii/faust/pkg/service/utils"
	qh "github.com/qiniu/go-sdk/v7/storagev2/http_client"
	qu "github.com/qiniu/go-sdk/v7/storagev2/uploader"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	_ "golang.org/x/image/webp"
)

const (
	defaultExpires    = 3600
	fmtKodoReturnBody = `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`
)

type (
	ImageUploadResponse struct {
		Bucket string
		Key    string
		Size   int64
		URLs   []string
	}
	CertsUploadResponse struct {
		CommonName string
		Expiration time.Time
	}
)

type ServiceInterface interface {
	UploadImage(ctx context.Context, image string) (*ImageUploadResponse, error)
	UploadCerts(ctx context.Context, key, cert string) (*CertsUploadResponse, error)
	DeleteImage(ctx context.Context, name string) error
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
		bucketManager *storage.BucketManager
		uploader      *qu.UploadManager
	}
)

func NewQServiceConfig() *QServiceConfig {
	return &QServiceConfig{
		Expires: defaultExpires,
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
	mac := auth.New(cfg.AccessKey, cfg.SecretKey)
	return &QiniuService{
		config:        cfg,
		credentials:   mac,
		bucketManager: storage.NewBucketManager(mac, &storage.Config{}),
		uploader: qu.NewUploadManager(&qu.UploadManagerOptions{
			Options: qh.Options{
				Credentials: mac,
			},
		}),
	}
}

// UploadImage reads file and put on to the specific bucket.
func (s *QiniuService) UploadImage(ctx context.Context, name string) (*ImageUploadResponse, error) {
	format, size, err := image.DiscoverImage(name)
	if err != nil {
		return nil, err
	}

	// doc: https://developer.qiniu.com/kodo/1238/go#upload-file
	key := utils.GenUploadKey(format)
	if err := s.uploader.UploadFile(ctx, name, &qu.ObjectOptions{
		BucketName:  s.config.Bucket,
		ObjectName:  &key,
		FileName:    key,
		ContentType: "application/json",
	}, nil); err != nil {
		return nil, err
	}

	// query bucket domains
	domainInfo, err := s.bucketManager.ListBucketDomains(s.config.Bucket)
	if err != nil {
		return nil, err
	}
	urls := make([]string, len(domainInfo))
	for i, domain := range domainInfo {
		urls[i] = fmt.Sprintf("https://%s/%s", domain.Domain, key)
	}

	return &ImageUploadResponse{
		Bucket: s.config.Bucket,
		Key:    key,
		Size:   size,
		URLs:   urls,
	}, nil
}

func (s *QiniuService) DeleteImage(ctx context.Context, name string) error {
	return s.bucketManager.Delete(s.config.Bucket, name)
}

type SSLCerts struct {
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
	Key        string `json:"pri"`
	CertChain  string `json:"ca"`
}

func (s *QiniuService) UploadCerts(ctx context.Context, keyPath, certPath string) (*CertsUploadResponse, error) {
	rawKeyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key %s: %w", keyPath, err)
	}
	block, _ := pem.Decode(rawKeyData)
	if block == nil {
		return nil, errors.New("invalid private key: no PEM data found")
	}

	rawCertData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificates %s: %w", certPath, err)
	}

	// Decode and parse all PEM blocks to support full certificate chains
	var cert *x509.Certificate
	rest := rawCertData
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		parsedCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %w", err)
		}
		if cert == nil {
			cert = parsedCert // leaf certificate (first in chain)
		}
	}
	if cert == nil {
		return nil, errors.New("invalid certificate: no PEM data found")
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
	if _, err := postRequest(ctx, s.credentials, "/sslcert", reqBody); err != nil {
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
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ResponseError(resp)
	}
	return io.ReadAll(resp.Body)
}
