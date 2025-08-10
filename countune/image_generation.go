package countune

import (
	// local
	"countube/common"

	// standard
	"fmt"
	"image"
	"image/color"

	"image/draw"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	// external
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

type countuneMeta struct {
	id     int
	barSeq int
	bars   int
}

func PrepareImagesForVideo(vidCfg VideoConfig) {
	common.EnsurePath(OutputPath)

	numCountuneBars, numCountuneBarsForLoop, countuneWidth := calculateCountuneSizeForVideo(vidCfg)
	countuneVideoImg := image.NewRGBA(image.Rect(0, 0, countuneWidth, vidCfg.VideoHeight))
	bgImage := image.NewUniform(vidCfg.BackgroundColor)
	pasteImage(bgImage, countuneVideoImg, 0, 0)

	fmt.Println("Generating title screen...")
	generateVideoTitle(vidCfg)

	fmt.Printf("Generating random Countune picture with %d bars...\n", numCountuneBars)
	countuneImg := generateRandomCountuneForVideo(vidCfg.Name, numCountuneBars, numCountuneBarsForLoop)

	fmt.Println("Resizing Countune picture...")
	resizedCountuneImg := resize.Resize(0, uint(vidCfg.CountuneHeight), countuneImg, resize.NearestNeighbor)

	// we center the countune stripe vertically on the center of the video screen
	y := (vidCfg.VideoHeight - resizedCountuneImg.Bounds().Dy()) / 2
	pasteImage(resizedCountuneImg, countuneVideoImg, 0, y)

	fmt.Println("Adding texts to video image...")
	drawTextOnVideoImage(countuneVideoImg, vidCfg)

	outFilename := vidCfg.Name + OutputFullVideoImageFilenameExt
	common.WritePngToFile(filepath.Join(OutputPath, outFilename), countuneVideoImg)
	fmt.Println("Generated image for full video: ", outFilename)
}

func pasteImage(sourceImg image.Image, targetImg *image.RGBA, targetImgX int, targetImgY int) {
	srcImgBounds := sourceImg.Bounds()
	tgtImgBounds := image.Rectangle{
		image.Point{targetImgX, targetImgY},
		image.Point{targetImgX + srcImgBounds.Dx(), targetImgY + srcImgBounds.Dy()}}
	draw.Draw(targetImg, tgtImgBounds, sourceImg, image.ZP, draw.Src)
}

func generateVideoTitle(vidCfg VideoConfig) *image.RGBA {

	titlePic := generateVideoTitlePic(vidCfg.VideoWidth, vidCfg.VideoHeight, vidCfg.BackgroundColor,
		vidCfg.TitleUpperText, vidCfg.TitleLowerText)
	outFilename := vidCfg.Name + OutputTitlePicFilenameExt
	common.WritePngToFile(filepath.Join(OutputPath, outFilename), titlePic)

	return titlePic
}

func generateRandomCountuneForVideo(videoName string, numBars int, numBarsToLoop int) *image.RGBA {

	mainPicMetas := selectRandomCountunePics(numBars)
	loopPicMetas := collectNCountuneBars(mainPicMetas, numBarsToLoop)
	picMetas := append(mainPicMetas, loopPicMetas...)

	img := combineCountunePics(picMetas)
	outFilename := videoName + OutputRandomCountuneFilenameExt
	common.WritePngToFile(filepath.Join(OutputPath, outFilename), img)

	return img
}

func collectNCountuneBars(picMetas []countuneMeta, n int) []countuneMeta {
	var result []countuneMeta
	total := 0

	for _, elem := range picMetas {
		if total+elem.bars <= n {
			result = append(result, elem)
			total += elem.bars
		} else {
			remaining := n - total
			if remaining > 0 {
				partial := elem
				partial.bars = remaining
				result = append(result, partial)
			}
			break
		}
	}
	return result
}

func calculateCountuneSizeForVideo(vidCfg VideoConfig) (int, int, int) {
	outputBarWidth := vidCfg.CountuneHeight / CountunePicBarWidthToHeightRatio
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

func selectRandomCountunePics(totalBars int) []countuneMeta {

	availablePics := scanCountunePics(CountunePicCachePath)
	fmt.Printf("Found %d Countune pics to choose from.\n", len(availablePics))

	selectedPics := make([]countuneMeta, 0, 100)
	rand.Seed(time.Now().UnixNano())

	for remainingBars := totalBars; remainingBars > 0; {

		if len(availablePics) == 0 {
			log.Fatal("Not enough pictures available for a video of such length.")
		}

		// select a random picture from the ones available and remove it from the list
		selectedPicId := rand.Intn(len(availablePics))
		selectedPic := availablePics[selectedPicId]

		// TODO: use a bit set
		// remove the selected picture from the list of availables
		tmp := make([]countuneMeta, 0, cap(availablePics))
		tmp = append(tmp, availablePics[:selectedPicId]...)
		tmp = append(tmp, availablePics[selectedPicId+1:]...)
		availablePics = tmp

		// add the selected pic to the results
		effectiveBars := common.Min(remainingBars, selectedPic.bars)
		selectedPic.bars = effectiveBars
		selectedPics = append(selectedPics, selectedPic)

		remainingBars -= effectiveBars
	}

	return selectedPics
}

func combineCountunePics(pics []countuneMeta) *image.RGBA {
	totalBars := 0
	for i := 0; i < len(pics); i++ {
		totalBars += pics[i].bars
	}

	totalWidth := totalBars * CountunePicOriginalBarWidth
	outImageRect := image.Rectangle{image.Point{0, 0}, image.Point{totalWidth, CountunePicOriginalHeight}}
	outImage := image.NewRGBA(outImageRect)

	for i, x := 0, 0; i < len(pics); i++ {

		currentPic := pics[i]
		picFilePath := getPicFilePath(currentPic.id)
		picImage := common.ReadImageFromFile(picFilePath)

		picWidth := currentPic.bars * CountunePicOriginalBarWidth
		outBounds := image.Rectangle{image.Point{x, 0}, image.Point{x + picWidth, CountunePicOriginalHeight}}
		draw.Draw(outImage, outBounds, picImage, image.Point{0, 0}, draw.Src)

		x += picWidth
	}

	return outImage
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
	fontPath := FontPath

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

	img1 := drawTextOnImage(rgba, image.Rectangle{image.ZP, image.Point{width, height}}, float64(upperTextY), fontPath, upperTextFontSize, upperTextColor, upperTextString)
	rgba2 := image.NewRGBA(rgbaBounds)
	draw.Draw(rgba2, image.Rectangle{image.ZP, image.Point{width, height}}, img1, image.ZP, draw.Src)

	img2 := drawTextOnImage(rgba2, image.Rectangle{image.ZP, image.Point{width, height}}, float64(lowerTextY), fontPath, lowerTextFontSize, lowerTextColor, lowerTextString)
	rgba3 := image.NewRGBA(rgbaBounds)
	draw.Draw(rgba3, image.Rectangle{image.ZP, image.Point{width, height}}, img2, image.ZP, draw.Src)

	return rgba3
}

func drawTextOnImage(rgba *image.RGBA, bounds image.Rectangle, y float64, fontFilePath string, fontSize float64,
	color color.RGBA, text string) image.Image {

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	dc := gg.NewContext(width, height)
	dc.DrawImage(rgba, 0, 0)

	err := dc.LoadFontFace(fontFilePath, fontSize)
	common.CheckErr(err)

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
	dc.LoadFontFace(FontPath, fontSize)
	dc.SetColor(color)
	vidTxts := vidCfg.Texts
	barWidth := vidCfg.CountuneHeight / CountunePicBarWidthToHeightRatio

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
