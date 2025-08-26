package util

import (
	"errors"
)

type downsampler struct {
	framesToKeepPerSample        int
	framesToDropPerSample        int
	sampleSize                   int
	frameIdxWithinSample         int
	framesKeptInCurrentSample    int
	framesDroppedInCurrentSample int
}

func NewDownsampler(sourceFps int, targetFps int) (error, *downsampler) {
	if targetFps <= 0 || targetFps > 60 {
		return errors.New("target fps must be a positive value lower than or equal to 60"), nil
	}
	if targetFps > sourceFps {
		return errors.New("target fps must be lower than or equal to the source fps."), nil
	}

	// we don't really need to use the gcd with this algorithm but it's nice for debugging
	g := gcd(sourceFps, targetFps)
	framesToKeepPerSample := targetFps / g
	sampleSize := sourceFps / g

	return nil, &downsampler{
		framesToKeepPerSample:        framesToKeepPerSample,
		framesToDropPerSample:        sampleSize - framesToKeepPerSample,
		sampleSize:                   sampleSize,
		frameIdxWithinSample:         0,
		framesKeptInCurrentSample:    0,
		framesDroppedInCurrentSample: 0,
	}
}

func (ds *downsampler) ShouldDropNextFrame() bool {

	keepToDropRatioComp := ds.framesKeptInCurrentSample*ds.framesToDropPerSample - ds.framesDroppedInCurrentSample*ds.framesToKeepPerSample
	var drop bool

	if keepToDropRatioComp <= 0 {
		drop = false
		ds.framesKeptInCurrentSample++
	} else {
		drop = true
		ds.framesDroppedInCurrentSample++
	}

	ds.frameIdxWithinSample++

	if ds.frameIdxWithinSample == ds.sampleSize {
		ds.frameIdxWithinSample = 0
		ds.framesKeptInCurrentSample = 0
		ds.framesDroppedInCurrentSample = 0
	}

	return drop
}

/*
 * Euclidean algorithm for finding the Greatest Common Denominator (GCD)
 */
func gcd(a int, b int) int {

	if a == 0 {
		return b
	}

	return gcd(b%a, a)
}
