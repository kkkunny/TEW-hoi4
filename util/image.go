package util

import (
	"image"
	"os"

	"golang.org/x/image/draw"

	"github.com/ftrvxmtrx/tga"
)

func ResizeAndCopyTgaImage(f, t string, w, h uint16) error {
	inputFile, err := os.Open(f)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	fromImg, err := tga.Decode(inputFile)
	if err != nil {
		return err
	}

	toImg := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	draw.CatmullRom.Scale(toImg, toImg.Bounds(), fromImg, fromImg.Bounds(), draw.Over, nil)

	outputFile, err := os.Create(t)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	return tga.Encode(outputFile, toImg)
}
