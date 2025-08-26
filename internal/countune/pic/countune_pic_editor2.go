package pic

import (
	// "countube/internal/common"

	"image"
	"image/color"
	"image/draw"
)

// type BarPlacement int

// const (
// 	BELOW_WAVE BarPlacement = iota
// 	ABOVE_WAVE
// )

// func (bp BarPlacement) Toggle() BarPlacement {
// 	return 1 - bp
// }

type StripePlacement int

const (
	TOP StripePlacement = iota
	CENTER
	BOTTOM
)

type picEditor struct {
	backgroundColor color.Color
	barColor1       color.Color
	barColor2       color.Color
	canvas          *image.RGBA
	pixelsPerUnit   int
	// stripePlacement StripePlacement
	maxBarHeightUnits int
	topCoord          int
	barPlacement      BarPlacement
	barIdx            int
}

func NewPicEditor(picHeight int, pixelsPerUnit int, maxBars int, maxBarHeightUnits int, stripePlacement StripePlacement,
	initialBarPlacement BarPlacement) *picEditor {

	picWidth := maxBars * pixelsPerUnit
	canvas := image.NewRGBA(image.Rect(0, 0, picWidth, picHeight))
	// stripeTopCoord := 0

	return &picEditor{
		backgroundColor:   color.Black,
		barColor1:         color.Black,
		barColor2:         color.Black,
		canvas:            canvas,
		pixelsPerUnit:     pixelsPerUnit,
		maxBarHeightUnits: maxBarHeightUnits,
		//stripePlacement: stripePlacement,
		topCoord:     0, // fix this taking into account the stripe placement
		barPlacement: initialBarPlacement,
		barIdx:       0,
	}
}

func (e *picEditor) image() *image.RGBA {
	return e.canvas
}

func (e *picEditor) ChangeColors(backgroundColor color.Color, barColor1 color.Color, barColor2 color.Color) {
	e.backgroundColor = backgroundColor
	e.barColor1 = barColor1
	e.barColor2 = barColor2
}

func (e *picEditor) ToggleBarPlacement() {
	e.barPlacement.Toggle()
}

func (e *picEditor) DrawBar2(waveHeightInUnits int, barPlacement BarPlacement) {
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
