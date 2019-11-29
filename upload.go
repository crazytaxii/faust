package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	fp "path/filepath"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/urfave/cli"
)

type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func (fc *FaustClient) Upload(c *cli.Context) error {
	putPolicy := storage.PutPolicy{
		Scope:      fc.Config.Bucket,
		Expires:    3600,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	mac := qbox.NewMac(fc.Config.AccessKey, fc.Config.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := &storage.Config{}
	formUploader := storage.NewFormUploader(cfg)
	ret := &MyPutRet{}

	putExtra := &storage.PutExtra{}
	key := fmt.Sprintf("%s/%s.jpg", getDate(), randNum(10000000, 100000000))

	srcImg, err := os.Open(fc.ImgPath)
	if err != nil {
		return err
	}
	defer srcImg.Close()

	var outBuf []byte
	switch fp.Ext(fc.ImgPath) {
	case "." + getImgExtension(JPG):
		outBuf, err = ioutil.ReadAll(srcImg)
		if err != nil {
			return err
		}
	case "." + getImgExtension(PNG):
		fallthrough
	case "." + getImgExtension(GIF):
		fallthrough
	case "." + getImgExtension(BMP):
		fallthrough
	case "." + getImgExtension(WEBP):
		fallthrough
	case "." + getImgExtension(TIFF):
		outBuf, err = imgConvert(fc.ImgPath, JPG)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported image type")
	}

	err = formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(outBuf), int64(len(outBuf)), putExtra)
	if err != nil {
		return err
	}
	fmt.Println("bucket:", ret.Bucket)
	fmt.Println("key:", ret.Key)
	fmt.Println("file size:", ret.Fsize)
	fmt.Println("hash:", ret.Hash)
	fmt.Println("public access url:", storage.MakePublicURL(fc.Config.BaseUrl, ret.Key))
	return nil
}
