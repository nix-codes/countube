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
	"math/rand"
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
	mainSpec             CountuneSpec
	getPicSpec           func() *CountunePicSpec
	x                    int
	waveStep             int
	shouldChangeWaveStep bool
	barPlacement         BarPlacement
}

func NewOnTheFlyPicSeq(mainSpec CountuneSpec, getPicSpec func() *CountunePicSpec) *OnTheFlyPicSeq {

	return &OnTheFlyPicSeq{
		mainSpec:             mainSpec,
		getPicSpec:           getPicSpec,
		x:                    mainSpec.InitialNum,
		waveStep:             mainSpec.WaveStep,
		shouldChangeWaveStep: false,
		barPlacement:         mainSpec.InitialBarPlacement,
	}
}

func (s *OnTheFlyPicSeq) Next() image.Image {
	picSpec := s.getPicSpec()
	if picSpec == nil {
		return nil
	}

	if !s.shouldChangeWaveStep {
		s.shouldChangeWaveStep = common.Chance(s.mainSpec.WaveStepChangeProb)
	}

	img, newBarPlacement := s.drawPic(*picSpec)
	s.x += picSpec.NumBars
	s.barPlacement = newBarPlacement

	return img
}

func (s *OnTheFlyPicSeq) drawPic(picSpec CountunePicSpec) (*image.RGBA, BarPlacement) {

	editor := newPicEditor(
		s.mainSpec.PicHeight,
		s.mainSpec.BarWidth,
		picSpec.NumBars,
		s.mainSpec.WaveAmplitude,
		CENTER)
	editor.changeColors(
		hexToColor(picSpec.BackgroundColor),
		hexToColor(picSpec.BarColor1),
		hexToColor(picSpec.BarColor2))

	x := s.x
	nextX := x + picSpec.NumBars
	yPrev := s.waveFn(max(0, x-1))
	for ; x < nextX; x += 1 {

		if common.IsPrime(x) {
			s.barPlacement = s.barPlacement.Toggle()
		}

		y := s.waveFn(x)

		if s.shouldChangeWaveStep && y != yPrev {
			newStep := rand.Intn(s.mainSpec.WaveStep) + 1
			s.waveStep = newStep
			s.shouldChangeWaveStep = false

			waveDir := DOWN
			if y > yPrev {
				waveDir = UP
			}

			s.mainSpec.InitialWaveHeight = y
			s.mainSpec.InitialWaveDirection = waveDir
			s.mainSpec.InitialNum = x
		}

		yPrev = y

		editor.drawBar(y, s.barPlacement)
	}

	return editor.image(), s.barPlacement
}

func (s *OnTheFlyPicSeq) waveFn(x int) int {

	return countuneFn(
		x,
		s.waveStep,
		s.mainSpec.WaveAmplitude,
		s.mainSpec.InitialNum,
		s.mainSpec.InitialWaveHeight,
		int(s.mainSpec.InitialWaveDirection))

}

func hexToColor(hex string) color.Color {
	var r, g, b uint8
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
