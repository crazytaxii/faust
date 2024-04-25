package avif

import (
	"image"
	"io"
)

const avifHeader = "\x00\x00\x00?ftypavif"

func Decode(r io.Reader) (img image.Image, err error) {
	// It's unnecessary to implement this function
	return
}

func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	// It's unnecessary to implement this function
	return
}

func init() {
	image.RegisterFormat("avif", avifHeader, Decode, DecodeConfig)
}
