package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

type FileType int

const (
	PNG FileType = iota
	JPG
	GIF
	WEBP
	BMP
	TIFF
	ERR
)

func getDate() string {
	return time.Now().Format("06-01-02")
}

func randNum(min, max int) string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%d", rand.Intn(max-min)+min)
}

func getImgType(input string) FileType {
	switch input {
	case "jpg":
		fallthrough
	case "jpeg":
		return JPG
	case "png":
		return PNG
	case "gif":
		return GIF
	case "bmp":
		return BMP
	case "webp":
		return WEBP
	case "tiff":
		return TIFF
	default:
		return ERR
	}
}

func getImgExtension(input FileType) string {
	switch input {
	case JPG:
		return "jpg"
	case PNG:
		return "png"
	case GIF:
		return "gif"
	case BMP:
		return "bmp"
	case WEBP:
		return "webp"
	case TIFF:
		return "tiff"
	default:
		return ""
	}
}

func getImgName(srcPath string) string {
	ext := strings.ToLower(filepath.Ext(srcPath))
	_, fileName := filepath.Split(srcPath)
	return fileName[0 : len(fileName)-len(ext)]
}

func imgConvert(srcPath string, fileType FileType) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(srcPath))

	originType := getImgType(ext[1:])
	if originType == ERR {
		return nil, errors.New("invalid input file type")
	}

	// open file
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode
	imgData, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// encode in new type
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	switch fileType {
	case JPG:
		err = jpeg.Encode(w, imgData, &jpeg.Options{Quality: 100})
	case PNG:
		err = png.Encode(w, imgData)
	case WEBP:
		err = webp.Encode(w, imgData, nil)
	case GIF:
		err = gif.Encode(w, imgData, nil)
	case BMP:
		err = bmp.Encode(w, imgData)
	case TIFF:
		err = tiff.Encode(w, imgData, nil)
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
