package util

import (
	// local
	"countube/internal/common"

	// standard
	"image"
	"image/color"
	"image/draw"
)

type hImageCompound struct {
	numImages  int
	images     []image.Image
	imgMetas   []compoundImageMeta
	totalWidth int
	imgIdx     int
	bgImage    image.Image
}

type compoundImageMeta struct {
	fromX int
	toX   int // non-inclusive
}

func NewHorizontalImageCompound(images []image.Image, bgColor color.Color) *hImageCompound {

	n := len(images)
	absX := 0
	imgMetas := make([]compoundImageMeta, n)

	for i := 0; i < n; i++ {

		currentImgWidth := images[i].Bounds().Dx()
		meta := compoundImageMeta{
			fromX: absX,
			toX:   absX + currentImgWidth,
		}

		imgMetas[i] = meta
		absX = meta.toX
	}

	return &hImageCompound{
		numImages:  n,
		images:     images,
		imgMetas:   imgMetas,
		totalWidth: absX,
		imgIdx:     0,
		bgImage:    image.NewUniform(bgColor),
	}
}

func (ic *hImageCompound) Draw(targetImage *image.RGBA, fromX int, toX int) {

	ic.locateX(fromX)
	requestedWidth := toX - fromX
	requestedHeight := targetImage.Bounds().Dy()
	remainingWidth := requestedWidth
	tgtX := 0

	for i := ic.imgIdx; remainingWidth > 0 && i < ic.numImages; i++ {

		ic.imgIdx = i
		currImg := ic.images[ic.imgIdx]
		currImgMeta := ic.imgMetas[ic.imgIdx]
		widthToUseFromCurrImg := common.Min(remainingWidth, currImgMeta.toX-fromX)
		tgtBounds := image.Rect(tgtX, 0, tgtX+widthToUseFromCurrImg, requestedHeight)

		draw.Draw(targetImage, tgtBounds, currImg, image.Point{fromX - currImgMeta.fromX, 0}, draw.Src)

		remainingWidth = remainingWidth - widthToUseFromCurrImg
		fromX += widthToUseFromCurrImg
		tgtX += widthToUseFromCurrImg
	}

	if remainingWidth > 0 {
		// fill the remaining target area with black
		tgtBounds := image.Rect(tgtX, 0, tgtX+remainingWidth, requestedHeight)
		draw.Draw(targetImage, tgtBounds, ic.bgImage, image.ZP, draw.Src)
	}

}

func (ic *hImageCompound) locateX(x int) {

	if ic.imgIdx >= ic.numImages {
		ic.imgIdx = ic.numImages - 1
	}

	for ic.imgIdx < ic.numImages {
		meta := ic.imgMetas[ic.imgIdx]

		if x >= meta.fromX && x < meta.toX {
			break
		} else if x >= meta.fromX {
			ic.imgIdx++
		} else {
			ic.imgIdx--
		}
	}
}
