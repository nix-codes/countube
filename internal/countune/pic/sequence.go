package pic

/**
 ** Represents a sequence of Countune pictures, were the wave function
 ** is kept consistently across those pictures.
 **/

import (
	// local
	"countube/internal/common"

	// standard
	"fmt"
	"image"
	"image/color"
)

type FixedPicSeq struct {
	picSeq   *OnTheFlyPicSeq
	picSpecs []CountunePicSpec
	specIdx  int
}

func NewFixedPicSeq(mainSpec CountuneSpec, picSpecs []CountunePicSpec) *FixedPicSeq {

	seq := &FixedPicSeq{
		picSpecs: picSpecs,
		specIdx:  0,
	}

	getPicSpecFn := func() *CountunePicSpec {
		if seq.specIdx >= len(seq.picSpecs) {
			return nil
		}
		picSpec := &seq.picSpecs[seq.specIdx]
		seq.specIdx += 1
		return picSpec
	}

	seq.picSeq = NewOnTheFlyPicSeq(mainSpec, getPicSpecFn)

	return seq
}

func (s *FixedPicSeq) Next() image.Image {
	return s.picSeq.Next()
}

type OnTheFlyPicSeq struct {
	mainSpec     CountuneSpec
	getPicSpec   func() *CountunePicSpec
	x            int
	barPlacement BarPlacement
}

func NewOnTheFlyPicSeq(mainSpec CountuneSpec, getPicSpec func() *CountunePicSpec) *OnTheFlyPicSeq {

	return &OnTheFlyPicSeq{
		mainSpec:     mainSpec,
		getPicSpec:   getPicSpec,
		x:            mainSpec.InitialNum,
		barPlacement: mainSpec.InitialBarPlacement,
	}
}

func (s *OnTheFlyPicSeq) Next() image.Image {
	picSpec := s.getPicSpec()
	if picSpec == nil {
		return nil
	}

	img, newBarPlacement := drawPic(s.mainSpec, *picSpec, s.x, s.barPlacement)
	s.x += picSpec.NumBars
	s.barPlacement = newBarPlacement

	return img
}

func drawPic(mainSpec CountuneSpec, picSpec CountunePicSpec, x int, barPlacement BarPlacement) (*image.RGBA, BarPlacement) {

	editor := newPicEditor(
		mainSpec.PicHeight,
		mainSpec.BarWidth,
		picSpec.NumBars,
		mainSpec.WaveAmplitude,
		CENTER)
	editor.changeColors(
		hexToColor(picSpec.BackgroundColor),
		hexToColor(picSpec.BarColor1),
		hexToColor(picSpec.BarColor2))

	nextX := x + picSpec.NumBars
	for ; x < nextX; x += 1 {

		if common.IsPrime(x) {
			barPlacement = barPlacement.Toggle()
		}

		y := countuneFn(
			x,
			mainSpec.WaveStep,
			mainSpec.WaveAmplitude,
			mainSpec.InitialNum,
			mainSpec.InitialWaveHeight,
			int(mainSpec.InitialWaveDirection))

		editor.drawBar(y, barPlacement)
	}

	return editor.image(), barPlacement
}

func hexToColor(hex string) color.Color {
	var r, g, b uint8
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
