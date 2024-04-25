package image

import (
	"image"
	"io"
)

const avifHeader = "\x00\x00\x00?ftypavif"
const svgHeader = "<svg width=\""

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
