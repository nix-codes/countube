package video

import (
	// local
	"countube/assets"
	"countube/internal/api"
	"countube/internal/common"
	"countube/internal/countune"
	"countube/internal/countune/pic"

	// standard
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"path/filepath"
	"strings"

	// external
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func PrepareImagesForVideo(vidCfg api.VideoConfig) {
	common.EnsurePath(api.OutputPath)

	numCountuneBars, requiredBarsForLoop, countuneWidth := calculateCountuneSizeForVideo(vidCfg)
	countuneVideoImg := image.NewRGBA(image.Rect(0, 0, countuneWidth, vidCfg.VideoHeight))
	bgImage := image.NewUniform(vidCfg.BackgroundColor)
	pasteImage(bgImage, countuneVideoImg, 0, 0)

	generateVideoTitle(vidCfg)

	countuneImg := generateRandomCountuneForVideo(vidCfg.Name, numCountuneBars, requiredBarsForLoop)

	fmt.Println("Resizing Countune picture...")
	resizedCountuneImg := resize.Resize(0, uint(vidCfg.CountuneHeight), countuneImg, resize.NearestNeighbor)

	// we center the countune stripe vertically on the center of the video screen
	y := (vidCfg.VideoHeight - resizedCountuneImg.Bounds().Dy()) / 2
	pasteImage(resizedCountuneImg, countuneVideoImg, 0, y)

	fmt.Println("Adding texts to video image...")
	drawTextOnVideoImage(countuneVideoImg, vidCfg)

	outFilename := vidCfg.Name + api.OutputFullVideoImageFilenameExt
	common.WritePngToFile(filepath.Join(api.OutputPath, outFilename), countuneVideoImg)
	fmt.Println("Generated image for full video: ", outFilename)
}

func pasteImage(sourceImg image.Image, targetImg *image.RGBA, targetImgX int, targetImgY int) {
	srcImgBounds := sourceImg.Bounds()
	tgtImgBounds := image.Rectangle{
		image.Point{targetImgX, targetImgY},
		image.Point{targetImgX + srcImgBounds.Dx(), targetImgY + srcImgBounds.Dy()}}
	draw.Draw(targetImg, tgtImgBounds, sourceImg, image.ZP, draw.Src)
}

func generateVideoTitle(vidCfg api.VideoConfig) *image.RGBA {
	fmt.Println("Generating title screen...")
	titlePic := generateVideoTitlePic(vidCfg.VideoWidth, vidCfg.VideoHeight, vidCfg.BackgroundColor,
		vidCfg.TitleUpperText, vidCfg.TitleLowerText)
	outFilename := vidCfg.Name + api.OutputTitlePicFilenameExt
	common.WritePngToFile(filepath.Join(api.OutputPath, outFilename), titlePic)
	fmt.Println("Wrote " + outFilename)
	fmt.Println()

	return titlePic
}

func generateRandomCountuneForVideo(videoName string, numBars int, numBarsToLoop int) *image.RGBA {

	fmt.Printf("Generating random Countune picture with %d bars...\n", numBars)
	img := pic.GenerateRandomCountune(numBars)
	img = appendImageLeftSegment(img, numBarsToLoop*countune.CountunePicOriginalBarWidth)

	outFilename := videoName + api.OutputRandomCountuneFilenameExt
	common.WritePngToFile(filepath.Join(api.OutputPath, outFilename), img)
	fmt.Println("Wrote " + outFilename)
	fmt.Println()

	return img
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

func calculateCountuneSizeForVideo(vidCfg api.VideoConfig) (int, int, int) {
	outputBarWidth := vidCfg.CountuneHeight / countune.CountunePicBarWidthToHeightRatio
	singleScreenScrollSeconds := float64(vidCfg.VideoWidth) / float64(outputBarWidth) / vidCfg.CountuneSpeed
	titleScrollSeconds := float64(vidCfg.TitleDelay) + float64(singleScreenScrollSeconds)
	requiredBarsForLoop := 0

	if vidCfg.Loop {
		requiredBarsForLoop = int(singleScreenScrollSeconds * vidCfg.CountuneSpeed)
		titleScrollSeconds = 0
	}

	countuneScrollSeconds := float64(vidCfg.VideoLen) - titleScrollSeconds
	requiredBarsForVideo := int(countuneScrollSeconds * vidCfg.CountuneSpeed)
	countuneWidth := (requiredBarsForVideo + requiredBarsForLoop) * outputBarWidth

	return requiredBarsForVideo, requiredBarsForLoop, countuneWidth
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

func drawTextOnVideoImage(img *image.RGBA, vidCfg api.VideoConfig) image.Image {

	fontSize := 50.0
	color := color.RGBA{230, 150, 60, 255}

	dc := gg.NewContextForRGBA(img)
	dc.SetFontFace(assets.WhitrabtFont(fontSize))

	dc.SetColor(color)
	vidTxts := vidCfg.Texts
	barWidth := vidCfg.CountuneHeight / countune.CountunePicBarWidthToHeightRatio

	for i := 0; i < len(vidTxts); i++ {
		vidTxt := vidTxts[i]
		emptyUpperHeight := (vidCfg.VideoHeight - vidCfg.CountuneHeight) / 2
		_, h := dc.MeasureString(vidTxt.Text)
		y := (float64(emptyUpperHeight)-h)/2.0 + h

		startBar := common.Max(0, int(vidCfg.CountuneSpeed*float64(vidTxt.StartSeconds-vidCfg.TitleDelay)))
		x := startBar * barWidth
		dc.DrawString(vidTxt.Text, float64(x), y)
	}

	draw.Draw(img, img.Bounds(), dc.Image(), image.ZP, draw.Src)

	return dc.Image()
}
