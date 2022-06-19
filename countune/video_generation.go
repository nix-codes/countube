package countune

import (
	// local
	"countube/common"
	"countube/downsampler"

	// standard
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"time"
)

func GenerateImageScrollVideo(vidCfg VideoConfig) {

	titleImgFilename := vidCfg.Name + OutputTitlePicFilenameExt
	titleImgFilePath := filepath.Join(OutputPath, titleImgFilename)
	titleImg := common.ReadImageFromFile(titleImgFilePath)

	videoImgFilename := vidCfg.Name + OutputFullVideoImageFilenameExt
	videoImgFilePath := filepath.Join(OutputPath, videoImgFilename)
	videoImg := common.ReadImageFromFile(videoImgFilePath)

	framesFilename := vidCfg.Name + OutputVideoFramesFilenameExt
	outFile, err := os.Create(filepath.Join(OutputPath, framesFilename))
	common.CheckErr(err)
	defer outFile.Close()

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(outFile)

	resizedCountuneBarWidth := vidCfg.CountuneHeight / CountunePicBarWidthToHeightRatio
	pixelsPerSec := int(float64(resizedCountuneBarWidth) * vidCfg.CountuneSpeed)

	i := 0
	totalFrames := vidCfg.VideoLen * vidCfg.Fps
	startTime := time.Now()
	elapsed := time.Duration(0)
	lastElapsedCheckTime := startTime
	eta := elapsed

	generateImageScrollVideoFrames(titleImg, vidCfg.TitleDelay, videoImg, vidCfg.BackgroundColor,
		pixelsPerSec, vidCfg.Fps, func(frame image.Image) {

			fmt.Printf("processing frame: %d / %d  |  elapsed: %s  |  eta: %s                      \r",
				i+1, totalFrames, elapsed, eta)

			if time.Since(lastElapsedCheckTime).Truncate(time.Second) >= 1 {
				// we update the elapsed and eta every second
				lastElapsedCheckTime = time.Now()
				elapsed = time.Since(startTime).Truncate(time.Second)
				remainingFrames := totalFrames - i
				eta = time.Duration(int(float64(remainingFrames)*elapsed.Seconds()/float64(i))) * time.Second
			}

			// write frame into the file
			jpeg.Encode(buf, frame, &jpeg.Options{100})
			// common.WriteJpegToFile(fmt.Sprintf("out/frames/%04d.jpg", i), frame)

			_, err := writer.Write(buf.Bytes())
			common.CheckErr(err)
			buf.Reset()
			i++
		})
	writer.Flush()

	fmt.Println()
}

func generateImageScrollVideoFrames(titleImg image.Image, titleDelay int, mainImg image.Image, bgColor color.Color,
	pixelsPerSec int, fps int, frameProcessFn func(image.Image)) {

	videoWidth := titleImg.Bounds().Dx()
	videoHeight := titleImg.Bounds().Dy()

	totalSeconds := int(math.Ceil(float64(titleDelay) + float64(videoWidth+mainImg.Bounds().Dx())/float64(pixelsPerSec)))
	requiredFrames := totalSeconds * fps
	numTitleFrames := titleDelay * fps

	compoundImg := NewHorizontalImageCompound([]image.Image{titleImg, mainImg}, bgColor)
	frame := image.NewRGBA(image.Rect(0, 0, videoWidth, videoHeight))
	compoundImg.Draw(frame, 0, videoWidth)

	maxFps := pixelsPerSec
	err, downsampler := downsampler.New(maxFps, fps)
	common.CheckErr(err)

	fmt.Println("total seconds  : ", totalSeconds)
	fmt.Println("required frames: ", requiredFrames)

	for i, x := 1, 0; i <= requiredFrames; i++ {

		if i > numTitleFrames {

			// drop frames to fit target fps
			for downsampler.ShouldDropNextFrame() {
				x++
			}

			compoundImg.Draw(frame, x, x+videoWidth)
			x++
		}
		frameProcessFn(frame)
	}
}
