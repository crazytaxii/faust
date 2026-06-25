package image

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"image"
	"io"
	"os"
)

const avifHeader = "\x00\x00\x00?ftypavif"
const svgHeader = "<svg"

func decode(r io.Reader) (img image.Image, err error) {
	// It's unnecessary to implement this function
	return
}

func decodeConfig(r io.Reader) (cfg image.Config, err error) {
	// It's unnecessary to implement this function
	return
}

func init() {
	image.RegisterFormat("avif", avifHeader, decode, decodeConfig)
	image.RegisterFormat("svg", svgHeader, decode, decodeConfig)
}

// DiscoverImage detects the image format and returns the size.
func DiscoverImage(path string) (format string, size int64, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer func() { _ = f.Close() }()

	fi, err := f.Stat()
	if err != nil {
		return "", 0, err
	}
	_, format, err = image.DecodeConfig(f)
	if err != nil {
		return "", 0, err
	}
	return format, fi.Size(), nil
}
