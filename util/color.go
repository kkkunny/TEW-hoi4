package util

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

func NewColorByMode(mode string, v1, v2, v3 float64) (color.Color, error) {
	switch strings.ToLower(mode) {
	case "rgb", "":
		return NewRGB(uint8(v1), uint8(v2), uint8(v3)), nil
	case "hsv":
		return colorful.Hsv(v1, v2, v3), nil
	default:
		return nil, fmt.Errorf("unknown color type `%s`", mode)
	}
}

func NewRGB(r, g, b uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func GetRGBA(c color.Color) (uint8, uint8, uint8, uint8) {
	r, g, b, a := c.RGBA()
	return uint8(r & 0xFF), uint8(g & 0xFF), uint8(b & 0xFF), uint8(a & 0xFF)
}

func GetRGB(c color.Color) (uint8, uint8, uint8) {
	r, g, b, _ := GetRGBA(c)
	return r, g, b
}

func AlphaBlendColor(from, to color.Color, ratio float32) color.Color {
	ra, ga, ba, aa := from.RGBA()
	rb, gb, bb, _ := to.RGBA()
	fRatio := float32(aa) * ratio
	r := (float32(ra)*(65535-fRatio) + float32(rb)*fRatio) / 65535
	g := (float32(ga)*(65535-fRatio) + float32(gb)*fRatio) / 65535
	b := (float32(ba)*(65535-fRatio) + float32(bb)*fRatio) / 65535
	return NewRGB(uint8(r/257), uint8(g/257), uint8(b/257))
}
