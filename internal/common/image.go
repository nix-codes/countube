package common

import (
	// standard
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
)

func ConcatImages(leftImg image.Image, rightImg image.Image) *image.RGBA {

	lBounds := leftImg.Bounds()
	rBounds := rightImg.Bounds()
	maxHeight := Max(lBounds.Dy(), rBounds.Dy())

	outBounds := image.Rectangle{
		image.ZP,
		image.Point{lBounds.Dx() + rBounds.Dx(), maxHeight}}
	outRgba := image.NewRGBA(outBounds)
	ConcatImagesIntoTarget(leftImg, rightImg, outRgba)

	return outRgba
}

func ConcatImagesIntoTarget(leftImg image.Image, rightImg image.Image, targetImg *image.RGBA) {

	lBounds := leftImg.Bounds()
	rBounds := rightImg.Bounds()
	maxHeight := Max(lBounds.Dy(), rBounds.Dy())

	outYPos := (maxHeight - lBounds.Dy()) / 2
	outTmpBounds := image.Rectangle{
		image.Point{0, outYPos},
		image.Point{lBounds.Dx(), outYPos + lBounds.Dy()}}
	draw.Draw(targetImg, outTmpBounds, leftImg, image.ZP, draw.Src)

	outYPos = (maxHeight - rBounds.Dy()) / 2
	outTmpBounds = image.Rectangle{
		image.Point{lBounds.Dx(), outYPos},
		image.Point{lBounds.Dx() + rBounds.Dx(), outYPos + rBounds.Dy()}}

	draw.Draw(targetImg, outTmpBounds, rightImg, image.ZP, draw.Src)
}

func WritePngToFile(filePath string, img image.Image) {
	outFile, err := os.Create(filePath)
	CheckErr(err)
	defer outFile.Close()

	err = png.Encode(outFile, img)
	CheckErr(err)
}

func WriteJpegToFile(filePath string, img image.Image) {
	outFile, err := os.Create(filePath)
	CheckErr(err)
	defer outFile.Close()

	err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 100})
	CheckErr(err)
}

func ReadImageFromFile(filePath string) image.Image {

	f, err := os.Open(filePath)
	CheckErr(err)
	defer f.Close()

	img, _, err := image.Decode(f)
	CheckErr(err)

	return img
}

func StitchImagesHorizontally(images []image.Image) *image.RGBA {

	canvas := createStitchedImageCanvas(images)

	x := 0
	for _, img := range images {
		b := img.Bounds()
		dstRect := image.Rect(x, 0, x+b.Dx(), b.Dy())
		draw.Draw(canvas, dstRect, img, b.Min, draw.Over)
		x += b.Dx()
	}

	return canvas
}

func createStitchedImageCanvas(images []image.Image) *image.RGBA {
	totalWidth, maxHeight := calculateStitchedImageBounds(images)
	canvas := image.NewRGBA(image.Rect(0, 0, totalWidth, maxHeight))

	return canvas
}

func calculateStitchedImageBounds(images []image.Image) (int, int) {
	totalWidth := 0
	maxHeight := 0

	for _, img := range images {
		b := img.Bounds()
		totalWidth += b.Dx()
		if b.Dy() > maxHeight {
			maxHeight = b.Dy()
		}
	}

	return totalWidth, maxHeight
}
