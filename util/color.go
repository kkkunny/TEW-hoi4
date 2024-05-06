package util

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	stlbasic "github.com/kkkunny/stl/basic"
)

func Color(typeName string, v1, v2, v3 float64) (color.RGBA, error) {
	switch strings.ToLower(typeName) {
	case "rgb", "":
		return color.RGBA{R: uint8(v1), G: uint8(v2), B: uint8(v3)}, nil
	case "hsv":
		r, g, b := HSV2RGB(v1, v2, v3)
		return color.RGBA{R: r, G: g, B: b}, nil
	default:
		return stlbasic.Default[color.RGBA](), fmt.Errorf("unknown color type `%s`", typeName)
	}
}

func HSV2RGB(h, s, v float64) (uint8, uint8, uint8) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case 60 <= h && h < 120:
		r, g, b = x, c, 0
	case 120 <= h && h < 180:
		r, g, b = 0, c, x
	case 180 <= h && h < 240:
		r, g, b = 0, x, c
	case 240 <= h && h < 300:
		r, g, b = x, 0, c
	case 300 <= h:
		r, g, b = c, 0, x
	}
	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}
