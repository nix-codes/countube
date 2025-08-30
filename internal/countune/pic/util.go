package pic

/**
 ** Utility functions for analyzing Countune images.
 **/

import (
	"countube/internal/common"
	// "countube/internal/countune"

	"fmt"
	"image"
	"image/color"
)

func DetermineCountunePicSpec(img image.Image, barWidth int) (int, [3]string) {
	bounds := img.Bounds()
	width := bounds.Dx()
	barCount := width / barWidth
	colors := DetermineCountunePictureColors(img, barWidth)

	return barCount, colors
}

/**
 * Given an image with a Countune picture, it will determine the colors used in the image
 * (background color + the 2 colors for bars).
 *
 * The algorithm:
 * - Always pick 2 points when trying to get the color of a bar: one point on the top side and one on the bottom side.
 *   If the wave is not exactly at the bottom or the top, then you will get 2 different collors from those samples.
 *   One of those colors will correspond to the bar and the other to the background.
 * - We then take 2 more points for the next bar. This may suffice if again we get 2 different colors (we know that the
 *   color that repeats between this sample and the previous has to be the background color).
 * - If we weren't able to determine it (because the wave was at the bottom or the top and so there was only 1 color
 *   in the sample), we try to do the same looking at the 2 bars from the end of the image.
 *
 * Note: this algorithm only works for a wave with a period longer than the number of bars in the picture. If that
 *       condition is not met, the image could have the first two bars at the bottom/top and the last two at the
 *       bottom/top as well.
 */
func DetermineCountunePictureColors(img image.Image, barWidth int) [3]string {
	bounds := img.Bounds()
	width := bounds.Dx()
	barCount := width / barWidth

	colors0 := getColorSamplesFromBar(img, 0, barWidth)
	colors1 := getColorSamplesFromBar(img, 1, barWidth)
	colors, err := determineCountuneColors(colors0[0], colors0[1], colors1[0], colors1[1])

	if err != nil {
		colors0 = getColorSamplesFromBar(img, barCount-1, barWidth)
		colors1 = getColorSamplesFromBar(img, barCount-2, barWidth)
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
