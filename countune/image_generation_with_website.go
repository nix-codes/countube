package countune

import (
	// local
	"countube/common"

	// standard
	"fmt"
	"image"
	"image/draw"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	// external
	"github.com/nfnt/resize"
)

func PrepareImagesForVideo_old(vidCfg VideoConfig) {
	common.EnsurePath(OutputPath)

	numCountuneBars, numCountuneBarsForLoop, countuneWidth := calculateCountuneSizeForVideo(vidCfg)
	countuneVideoImg := image.NewRGBA(image.Rect(0, 0, countuneWidth, vidCfg.VideoHeight))
	bgImage := image.NewUniform(vidCfg.BackgroundColor)
	pasteImage(bgImage, countuneVideoImg, 0, 0)

	fmt.Println("Generating title screen...")
	generateVideoTitle(vidCfg)

	fmt.Printf("Generating random Countune picture with %d bars...\n", numCountuneBars)
	countuneImg := generateRandomCountuneForVideoUsingWebsitePics(vidCfg.Name, numCountuneBars, numCountuneBarsForLoop)

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

func generateRandomCountuneForVideoUsingWebsitePics(videoName string, numBars int, numBarsToLoop int) *image.RGBA {

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
