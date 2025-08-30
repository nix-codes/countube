package pic

import (
	"math/rand"
)

type BarPlacement int

const (
	BELOW_WAVE BarPlacement = iota
	ABOVE_WAVE
)

func (bp BarPlacement) Toggle() BarPlacement {
	return 1 - bp
}

type WaveDirection int

const (
	UP WaveDirection = iota
	DOWN
)

type CountuneSpec struct {
	PicHeight            int
	BarWidth             int
	WaveAmplitude        int
	WaveStep             int
	InitialNum           int
	InitialWaveHeight    int
	InitialWaveDirection WaveDirection
	InitialBarPlacement  BarPlacement
}

type CountunePicSpec struct {
	NumBars         int
	BarColor1       string
	BarColor2       string
	BackgroundColor string
}

var standardBarCounts = []int{
	16, 25, 36, 49, 64, 81,
}

var colorsInHex = []string{
	"800000", "f08080", "ff0000", "ff6347", "ffd700", "006400",
	"00ff00", "4169e1", "6495ed", "9400d3", "800080", "a9a9a9",
	"a0522d", "bc8f8f", "cd5c5c", "ff8c00", "ffff00", "008000",
	"808000", "0000ff", "000080", "ff00ff", "ffffff", "808080",
	"8b4513", "f4a460", "ff1493", "ffa500", "f0e68c", "9acd32",
	"7fff00", "87ceeb", "00ffff", "da70d6", "dcdcdc", "696969",
	"b22222", "deb887", "ff69b4", "ffa07a", "eee8aa", "adff2f",
	"8fbc8f", "b0c4de", "008080", "dda0dd", "d3d3d3", "708090",
	"cd853f", "d2b48c", "ffc0cb", "ffdab9", "fffacd", "bdb76b",
	"98fb98", "b0e0e6", "9370db", "d8bfd8", "c0c0c0", "000000",
}

const maxStartNum = 999999
const defaultWaveAmplitude = 10

func RandomCountuneSpec(picHeight, barWidth int) CountuneSpec {

	waveDirection := UP
	if rand.Intn(2) == 1 {
		waveDirection = DOWN
	}

	barPlacement := ABOVE_WAVE
	if rand.Intn(2) == 1 {
		barPlacement = BELOW_WAVE
	}

	return CountuneSpec{
		PicHeight:            picHeight,
		BarWidth:             barWidth,
		WaveAmplitude:        defaultWaveAmplitude,
		WaveStep:             5,
		InitialNum:           rand.Intn(maxStartNum + 1),
		InitialWaveHeight:    rand.Intn(defaultWaveAmplitude + 1),
		InitialWaveDirection: waveDirection,
		InitialBarPlacement:  barPlacement,
	}
}

func RandomPicSpec() CountunePicSpec {
	randomColors := selectRandomCountuneColors()

	return CountunePicSpec{
		NumBars:         selectRandomBarCount(),
		BackgroundColor: randomColors[0],
		BarColor1:       randomColors[1],
		BarColor2:       randomColors[2],
	}
}

func selectRandomBarCount() int {
	idx := rand.Intn(len(standardBarCounts))
	return standardBarCounts[idx]
}

func selectRandomCountuneColors() [3]string {
	nums := rand.Perm(60)
	bgColor := colorsInHex[nums[0]]
	barColor1 := colorsInHex[nums[1]]
	barColor2 := colorsInHex[nums[2]]

	return [3]string{bgColor, barColor1, barColor2}
}
