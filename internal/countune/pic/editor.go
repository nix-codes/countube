package pic

import (
	"image"
	"image/color"
	"image/draw"
)

type StripePlacement int

const (
	TOP StripePlacement = iota
	CENTER
	BOTTOM
)

type picEditor struct {
	backgroundColor   color.Color
	barColor1         color.Color
	barColor2         color.Color
	canvas            *image.RGBA
	pixelsPerUnit     int
	maxBarHeightUnits int
	topCoord          int
	barIdx            int
}

func newPicEditor(picHeight int, pixelsPerUnit int, maxBars int, maxBarHeightUnits int,
	stripePlacement StripePlacement) *picEditor {

	picWidth := maxBars * pixelsPerUnit
	canvas := image.NewRGBA(image.Rect(0, 0, picWidth, picHeight))

	stripeTopCoord := 0
	maxBarHeight := maxBarHeightUnits * pixelsPerUnit

	if stripePlacement == CENTER {
		stripeTopCoord = (picHeight - maxBarHeight) / 2
	} else if stripePlacement == BOTTOM {
		stripeTopCoord = picHeight - maxBarHeight
	}

	return &picEditor{
		backgroundColor:   color.Black,
		barColor1:         color.Black,
		barColor2:         color.Black,
		canvas:            canvas,
		pixelsPerUnit:     pixelsPerUnit,
		maxBarHeightUnits: maxBarHeightUnits,
		topCoord:          stripeTopCoord,
		barIdx:            0,
	}
}

func (e *picEditor) image() *image.RGBA {
	return e.canvas
}

func (e *picEditor) changeColors(backgroundColor color.Color, barColor1 color.Color, barColor2 color.Color) {
	e.backgroundColor = backgroundColor
	e.barColor1 = barColor1
	e.barColor2 = barColor2
}

func (e *picEditor) drawBar(waveHeightInUnits int, barPlacement BarPlacement) {
	var barColor color.Color

	if e.barIdx%2 == 0 {
		barColor = e.barColor1
	} else {
		barColor = e.barColor2
	}

	topColor := barColor
	bottomColor := e.backgroundColor

	if barPlacement == BELOW_WAVE {
		topColor = e.backgroundColor
		bottomColor = barColor
	}
	e.drawRectangle(0, e.maxBarHeightUnits-waveHeightInUnits, topColor)
	e.drawRectangle(e.maxBarHeightUnits-waveHeightInUnits, waveHeightInUnits, bottomColor)
	e.barIdx += 1
}

func (e *picEditor) drawRectangle(topUnits int, heightUnits int, color color.Color) {
	top := e.topCoord + topUnits*e.pixelsPerUnit
	bottom := top + heightUnits*e.pixelsPerUnit
	left := e.barIdx * e.pixelsPerUnit
	right := left + e.pixelsPerUnit

	bar := image.Rect(left, top, right, bottom)

	draw.Draw(e.canvas, bar, &image.Uniform{color}, image.Point{}, draw.Src)
}
