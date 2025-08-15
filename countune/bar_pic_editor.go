package countune

import (
	"countube/common"

	"image"
	"image/color"
	"image/draw"
)

type BarPlacement int

const (
	BELOW_WAVE BarPlacement = iota
	ABOVE_WAVE
)

func (bp BarPlacement) Toggle() BarPlacement {
	return 1 - bp
}

type barPicEditor struct {
	backgroundColor color.Color
	evenBarColor    color.Color
	oddBarColor     color.Color
	canvas          *image.RGBA
	barWidth        int
	unitHeight      int
	nextBarIdx      int
}

func NewBarPicEditor(picWidth int, picHeight int, barWidth int) *barPicEditor {

	canvas := image.NewRGBA(image.Rect(0, 0, picWidth, picHeight))

	return &barPicEditor{
		backgroundColor: color.Black,
		evenBarColor:    color.Black,
		oddBarColor:     color.Black,
		canvas:          canvas,
		barWidth:        barWidth,
		unitHeight:      int(picHeight / 10),
		nextBarIdx:      0,
	}
}

func (bpe *barPicEditor) image() *image.RGBA {
	return bpe.canvas
}

func (bpe *barPicEditor) ChangeColors(backgroundColor color.Color, evenBarColor color.Color, oddBarColor color.Color) {
	bpe.backgroundColor = backgroundColor
	bpe.evenBarColor = evenBarColor
	bpe.oddBarColor = oddBarColor
}

func (bpe *barPicEditor) DrawBar(waveAmplitude int, barPlacement BarPlacement) {
	var barColor color.Color

	if bpe.nextBarIdx%2 == 0 {
		barColor = bpe.evenBarColor
	} else {
		barColor = bpe.oddBarColor
	}

	x := bpe.nextBarIdx * bpe.barWidth

	if barPlacement == BELOW_WAVE {
		bpe.drawBarBelowWaveWithoutBackground(x, waveAmplitude, barColor)
		bpe.drawBarAboveWaveWithoutBackground(x, waveAmplitude, bpe.backgroundColor)
	} else {
		bpe.drawBarAboveWaveWithoutBackground(x, waveAmplitude, barColor)
		bpe.drawBarBelowWaveWithoutBackground(x, waveAmplitude, bpe.backgroundColor)
	}

	bpe.nextBarIdx += 1
}

func (bpe *barPicEditor) drawBarBelowWaveWithoutBackground(x int, waveAmplitude int, color color.Color) {
	bounds := bpe.canvas.Bounds()
	barHeight := waveAmplitude * bpe.unitHeight
	top := bounds.Max.Y - barHeight
	bottom := bounds.Max.Y

	bpe.drawBarWithoutBackground(x, top, bottom, color)
}

func (bpe *barPicEditor) drawBarAboveWaveWithoutBackground(x int, waveAmplitude int, color color.Color) {
	bounds := bpe.canvas.Bounds()
	barHeight := waveAmplitude * bpe.unitHeight
	top := 0
	bottom := bounds.Max.Y - barHeight

	bpe.drawBarWithoutBackground(x, top, bottom, color)
}

func (bpe *barPicEditor) drawBarWithoutBackground(x int, top int, bottom int, barColor color.Color) {
	bar := image.Rect(
		x,              // left
		top,            // top (Y grows downwards)
		x+bpe.barWidth, // right
		bottom,         // bottom
	)

	draw.Draw(bpe.canvas, bar, &image.Uniform{barColor}, image.Point{}, draw.Src)
}

func (bpe *barPicEditor) Write() {

	common.WritePngToFile("D:\\dev\\projects\\countube\\dummy.png", bpe.canvas)

}
