package pic

import (
	"countube/internal/common"
	"countube/internal/countune"

	"fmt"
	"image"
	"image/color"
)

func DetermineCountunePicSpec(img image.Image) (int, [3]string) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	barWidth := height / 10
	barCount := width / barWidth
	colors := DetermineCountunePictureColors(img)

	return barCount, colors
}

func DetermineCountunePictureColors(img image.Image) [3]string {
	bounds := img.Bounds()
	width := bounds.Dx()
	barCount := width / countune.CountunePicOriginalBarWidth

	colors0 := getColorSamplesFromBar(img, 0, countune.CountunePicOriginalBarWidth)
	colors1 := getColorSamplesFromBar(img, 1, countune.CountunePicOriginalBarWidth)
	colors, err := determineCountuneColors(colors0[0], colors0[1], colors1[0], colors1[1])

	if err != nil {
		colors0 = getColorSamplesFromBar(img, barCount-1, countune.CountunePicOriginalBarWidth)
		colors1 = getColorSamplesFromBar(img, barCount-2, countune.CountunePicOriginalBarWidth)
		colors, err = determineCountuneColors(colors0[0], colors0[1], colors1[0], colors1[1])
	}

	common.CheckErr(err)

	evenBarColor := colors[0]
	oddBarColor := colors[1]
	backgroundColor := colors[2]

	return [3]string{evenBarColor, oddBarColor, backgroundColor}
}

func determineCountuneColors(
	evenBarTopColor, evenBarBottomColor,
	oddBarTopColor, oddBarBottomColor string,
) ([]string, error) {

	// Count occurrences
	counts := make(map[string]int)
	colorsInput := []string{
		evenBarTopColor, evenBarBottomColor,
		oddBarTopColor, oddBarBottomColor,
	}
	for _, c := range colorsInput {
		counts[c]++
	}

	if len(counts) != 3 {
		return nil, fmt.Errorf("Cannot determine the colors effectively because we need 3 different colors in the samples.")
	}

	result := make([]string, 3)

	foundBackgroundColor := false
	for color, count := range counts {
		if count == 2 {
			result[2] = color
			foundBackgroundColor = true
			break
		}
	}

	if !foundBackgroundColor {
		return nil, fmt.Errorf("no color appears exactly 3 times")
	}

	if evenBarTopColor != result[2] {
		result[0] = evenBarTopColor
	} else {
		result[0] = evenBarBottomColor
	}

	if oddBarTopColor != result[2] {
		result[1] = oddBarTopColor
	} else {
		result[1] = oddBarBottomColor
	}

	return result, nil
}

func getColorSamplesFromBar(img image.Image, barIdx int, barWidth int) [2]string {
	height := img.Bounds().Dy()
	x := barIdx*barWidth + 1
	topColor := pixelColorHex(img, x, 1)
	bottomColor := pixelColorHex(img, x, height-1)

	return [2]string{topColor, bottomColor}
}

func pixelColorHex(img image.Image, x, y int) string {
	// we acocunt for the case where there's alpha channel
	c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
	return fmt.Sprintf("%02x%02x%02x", c.R, c.G, c.B)
}
