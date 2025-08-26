package pic

import (
	// local
	"countube/internal/common"

	// standard
	"fmt"
	"image"
	"image/color"
	"math/rand"
)

var standardBarCounts = []int{
	16, 25, 36, 49, 64, 81,
}

const barWidth = 20
const picHeight = 200

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

var colors []color.Color

type CountuneComposite struct {
	nextNum      int
	totalBars    int
	drawnBars    int
	barPlacement BarPlacement
	picEditor    *barPicEditor
}

func init() {
	colors = make([]color.Color, len(colorsInHex))
	for i, hex := range colorsInHex {
		colors[i] = hexToColor(hex)
	}
}

func GenerateCountuneFromSpec(spec *CountuneCompositeSpec) *image.RGBA {

	composite := NewCountuneComposite(spec.startNum, spec.totalBars)
	for _, picSpec := range spec.picSpecs {
		composite.addCountune(picSpec.barCount,
			picSpec.bgColor, picSpec.barColor1, picSpec.barColor2)
	}

	return composite.image()
}

func GenerateRandomCountune(barCount int) *image.RGBA {
	startNumber := rand.Intn(5000) + 1
	composite := NewCountuneComposite(startNumber, barCount)
	remainingBars := barCount

	for remainingBars > 0 {

		nextCountuneSize := selectRandomSize(remainingBars)
		randomColors := selectRandomCountuneColors()
		composite.addCountune(nextCountuneSize, randomColors[0], randomColors[1], randomColors[2])
		remainingBars -= nextCountuneSize
	}

	return composite.image()
}

func NewCountuneComposite(startNumber int, barCount int) *CountuneComposite {
	picWidth := barCount * barWidth
	picEditor := NewBarPicEditor(picWidth, picHeight, barWidth)

	primesBelowStartNumber := common.CountPrimesBelow(startNumber)
	var barPlacement BarPlacement
	if primesBelowStartNumber%2 == 0 {
		barPlacement = ABOVE_WAVE
	} else {
		barPlacement = BELOW_WAVE
	}

	return &CountuneComposite{
		nextNum:      startNumber,
		totalBars:    barCount,
		drawnBars:    0,
		barPlacement: barPlacement,
		picEditor:    picEditor,
	}
}

func (c *CountuneComposite) addCountune(barCount int, bgColor string, barColor1 string, barColor2 string) {

	c.picEditor.ChangeColors(hexToColor(bgColor), hexToColor(barColor1), hexToColor(barColor2))
	endNum := c.nextNum + barCount - 1

	for i := c.nextNum; i <= endNum; i += 1 {
		if common.IsPrime(i) {
			c.barPlacement = c.barPlacement.Toggle()
		}

		y := countuneFn(i, 5, 1)
		c.picEditor.DrawBar(y, c.barPlacement)
	}

	c.nextNum = endNum + 1
}

func (c *CountuneComposite) image() *image.RGBA {
	return c.picEditor.image()
}

func hexToColor(hex string) color.Color {
	var r, g, b uint8
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func selectRandomSize(remainingBarCount int) int {

	var filtered []int
	for _, val := range standardBarCounts {
		if val <= remainingBarCount {
			filtered = append(filtered, val)
		}
	}

	if len(filtered) == 0 {
		return remainingBarCount
	}

	return filtered[rand.Intn(len(filtered))]
}

func selectRandomPicSize() int {
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
