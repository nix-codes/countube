package generator

import (
	// local
	"countube/assets"
	"countube/internal/common"
	"countube/internal/countune/pic"

	// standard
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"path/filepath"
	"strings"

	// external
	"github.com/fogleman/gg"
)

const defaultWaveAmplitude = 10

func PrepareImagesForVideo(vidCfg VideoConfig) {
	common.EnsurePath(OutputPath)

	numCountuneBars, _ := calculateCountuneSizeForVideo(vidCfg)

	fmt.Println("Generating title screen...")
	titleImg := generateVideoTitlePic(vidCfg.ScreenWidth, vidCfg.ScreenHeight, vidCfg.BackgroundColor,
		vidCfg.TitleUpperText, vidCfg.TitleLowerText)

	fmt.Println("Generating picture sequence...")
	countuneStripImg := generateRandomCountuneForVideo(vidCfg, numCountuneBars, 0)

	fmt.Println("Adding texts to video image...")
	drawTextOnVideoImage(countuneStripImg, vidCfg)

	videoImg := common.StitchImagesHorizontally([]image.Image{titleImg, countuneStripImg})

	outFilename := vidCfg.Name + OutputFullVideoImageFilenameExt
	common.WritePngToFile(filepath.Join(OutputPath, outFilename), videoImg)
	fmt.Println("Generated image for full video: ", outFilename)
}

func pasteImage(sourceImg image.Image, targetImg *image.RGBA, targetImgX int, targetImgY int) {
	srcImgBounds := sourceImg.Bounds()
	tgtImgBounds := image.Rectangle{
		image.Point{targetImgX, targetImgY},
		image.Point{targetImgX + srcImgBounds.Dx(), targetImgY + srcImgBounds.Dy()}}
	draw.Draw(targetImg, tgtImgBounds, sourceImg, image.ZP, draw.Src)
}

func generateRandomCountuneForVideo(vidCfg VideoConfig, numBars, numBarsToLoop int) *image.RGBA {
	videoHeight := vidCfg.ScreenHeight
	countuneHeight := vidCfg.CountuneHeight

	fmt.Printf("Generating random Countune picture with %d bars...\n", numBars)
	img := generateRandomCountune(videoHeight, countuneHeight, numBars)

	return img
}

func generateRandomCountune(picHeight int, stripeHeight int, numBars int) *image.RGBA {
	if stripeHeight%defaultWaveAmplitude != 0 {
		log.Fatalf("Stripe height must be divisible by %d", defaultWaveAmplitude)
	}

	barWidth := stripeHeight / defaultWaveAmplitude

	countuneSpec := pic.RandomCountuneSpec(picHeight, barWidth)
	var picSpecs []pic.CountunePicSpec
	remainingBars := numBars

	for remainingBars > 0 {
		picSpec := pic.RandomPicSpec()
		adjustedNumBars := min(picSpec.NumBars, remainingBars)
		picSpec.NumBars = adjustedNumBars

		picSpecs = append(picSpecs, picSpec)
		remainingBars -= picSpec.NumBars
	}

	picSeq := pic.NewFixedPicSeq(countuneSpec, picSpecs)

	return pic.BuildCountuneStrip(picSeq)
}

func appendImageLeftSegment(img *image.RGBA, segmentWidth int) *image.RGBA {
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	height := bounds.Dy()

	newImg := image.NewRGBA(image.Rect(0, 0, origWidth+segmentWidth, height))

	// copy original image to new image
	draw.Draw(newImg, image.Rect(0, 0, origWidth, height), img, bounds.Min, draw.Src)

	// copy left segment to the appended area
	leftSegmentRect := image.Rect(0, 0, segmentWidth, height)
	draw.Draw(newImg, image.Rect(origWidth, 0, origWidth+segmentWidth, height), img, leftSegmentRect.Min, draw.Src)

	return newImg
}

func calculateCountuneSizeForVideo(vidCfg VideoConfig) (int, int) {
	scrollPixelsPerSec := vidCfg.ScrollSpeed * float64(vidCfg.BarWidth)
	singleScreenScrollSeconds := float64(vidCfg.ScreenWidth) / scrollPixelsPerSec
	titleScrollSeconds := float64(vidCfg.TitleDelay) + float64(singleScreenScrollSeconds)
	countuneScrollSeconds := float64(vidCfg.VideoLen) - titleScrollSeconds
	requiredBarsForVideo := int(countuneScrollSeconds * vidCfg.ScrollSpeed)
	countuneWidth := requiredBarsForVideo * vidCfg.BarWidth

	return requiredBarsForVideo, countuneWidth
}

// TODO: improve this mess
func generateVideoTitlePic(width int, height int, bgColor color.Color, upperText []string, lowerText []string) *image.RGBA {

	scaleFactor := float64(height) / 1080.0 // our default font size and spacing are based on a 1080 height screen
	f := 0.7                                // this is a hacky correction factor for the estimation of text boxes

	upperTextString := strings.Join(upperText, "\n")
	lowerTextString := strings.Join(lowerText, "\n")

	upperTextFontSize := 80.0 * scaleFactor
	lowerTextFontSize := 60.0 * scaleFactor
	upperTextColor := color.RGBA{0, 255, 0, 255}
	lowerTextColor := color.RGBA{0, 100, 255, 255}

	rgbaBounds := image.Rectangle{image.Point{0, 0}, image.Point{width, height}}
	rgba := image.NewRGBA(rgbaBounds)
	bgImg := image.NewUniform(bgColor)
	pasteImage(bgImg, rgba, 0, 0)

	numUpperTextLines := len(upperText)
	numLowerTextLines := len(lowerText)
	upperTextBoxSize := int(upperTextFontSize*f*float64(numUpperTextLines)) + int(upperTextFontSize*f*float64(numUpperTextLines-1)*1.2)
	lowerTextBoxSize := int(lowerTextFontSize*f*float64(numLowerTextLines)) + int(lowerTextFontSize*f*float64(numLowerTextLines-1)*1.2)
	spacingBoxSize := int(190.0 * scaleFactor) //int(upperTextFontSize * f * 2)
	totalTextHeight := upperTextBoxSize + spacingBoxSize + lowerTextBoxSize
	upperTextY := (height - totalTextHeight) / 2
	lowerTextY := upperTextY + totalTextHeight - lowerTextBoxSize

	img1 := drawTextOnImage(rgba, image.Rectangle{image.ZP, image.Point{width, height}}, float64(upperTextY), upperTextFontSize, upperTextColor, upperTextString)
	rgba2 := image.NewRGBA(rgbaBounds)
	draw.Draw(rgba2, image.Rectangle{image.ZP, image.Point{width, height}}, img1, image.ZP, draw.Src)

	img2 := drawTextOnImage(rgba2, image.Rectangle{image.ZP, image.Point{width, height}}, float64(lowerTextY), lowerTextFontSize, lowerTextColor, lowerTextString)
	rgba3 := image.NewRGBA(rgbaBounds)
	draw.Draw(rgba3, image.Rectangle{image.ZP, image.Point{width, height}}, img2, image.ZP, draw.Src)

	return rgba3
}

func drawTextOnImage(rgba *image.RGBA, bounds image.Rectangle, y float64, fontSize float64,
	color color.RGBA, text string) image.Image {

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	dc := gg.NewContext(width, height)
	dc.DrawImage(rgba, 0, 0)
	dc.SetFontFace(assets.WhitrabtFont(fontSize))

	x := float64(width / 2)
	maxWidth := float64(width) - float64(fontSize)
	dc.SetColor(color)
	dc.DrawStringWrapped(text, x, y, 0.5, 0, maxWidth, 1.8, gg.AlignCenter)

	return dc.Image()
}

func drawTextOnVideoImage(img *image.RGBA, vidCfg VideoConfig) image.Image {

	fontSize := 50.0
	color := color.RGBA{230, 150, 60, 255}

	dc := gg.NewContextForRGBA(img)
	dc.SetFontFace(assets.WhitrabtFont(fontSize))

	dc.SetColor(color)
	vidTxts := vidCfg.Texts

	for i := 0; i < len(vidTxts); i++ {
		vidTxt := vidTxts[i]
		emptyUpperHeight := (vidCfg.ScreenHeight - vidCfg.CountuneHeight) / 2
		_, h := dc.MeasureString(vidTxt.Text)
		y := (float64(emptyUpperHeight)-h)/2.0 + h

		startBar := common.Max(0, int(vidCfg.ScrollSpeed*float64(vidTxt.StartSeconds-vidCfg.TitleDelay)))
		x := startBar * vidCfg.BarWidth
		dc.DrawString(vidTxt.Text, float64(x), y)
	}

	draw.Draw(img, img.Bounds(), dc.Image(), image.ZP, draw.Src)

	return dc.Image()
}
