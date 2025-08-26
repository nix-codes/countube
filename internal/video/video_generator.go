package video

import (
	// local
	"countube/internal/api"
	"countube/internal/common"
	"countube/internal/countune"
	"countube/internal/video/util"

	// standard
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	"os"
	"path/filepath"
	"time"
)

func GenerateImageScrollVideo(vidCfg api.VideoConfig) {

	var titleImg image.Image

	if !vidCfg.Loop {
		titleImgFilename := vidCfg.Name + api.OutputTitlePicFilenameExt
		titleImgFilePath := filepath.Join(api.OutputPath, titleImgFilename)
		titleImg = common.ReadImageFromFile(titleImgFilePath)
	}

	videoImgFilename := vidCfg.Name + api.OutputFullVideoImageFilenameExt
	videoImgFilePath := filepath.Join(api.OutputPath, videoImgFilename)
	videoImg := common.ReadImageFromFile(videoImgFilePath)

	framesFilename := vidCfg.Name + api.OutputVideoFramesFilenameExt
	outFile, err := os.Create(filepath.Join(api.OutputPath, framesFilename))
	common.CheckErr(err)
	defer outFile.Close()

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(outFile)

	resizedCountuneBarWidth := vidCfg.CountuneHeight / countune.CountunePicBarWidthToHeightRatio
	pixelsPerSec := int(float64(resizedCountuneBarWidth) * vidCfg.CountuneSpeed)

	i := 0
	totalFrames := vidCfg.VideoLen * vidCfg.Fps
	startTime := time.Now()
	elapsed := time.Duration(0)
	lastElapsedCheckTime := startTime
	eta := elapsed

	generateImageScrollVideoFrames(vidCfg, titleImg, videoImg, pixelsPerSec,
		func(frame image.Image) {

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

func generateImageScrollVideoFrames(vidCfg api.VideoConfig, titleImg image.Image, mainImg image.Image,
	pixelsPerSec int, frameProcessFn func(image.Image)) {

	videoWidth := vidCfg.VideoWidth
	videoHeight := vidCfg.VideoHeight
	bgColor := vidCfg.BackgroundColor
	fps := vidCfg.Fps
	totalSeconds := vidCfg.VideoLen
	titleDelay := 0

	imagesToScroll := []image.Image{mainImg}

	if !vidCfg.Loop {
		imagesToScroll = append([]image.Image{titleImg}, imagesToScroll...)
		titleDelay = vidCfg.TitleDelay
	}

	requiredFrames := totalSeconds * fps
	numTitleFrames := titleDelay * fps

	compoundImg := util.NewHorizontalImageCompound(imagesToScroll, bgColor)
	frame := image.NewRGBA(image.Rect(0, 0, videoWidth, videoHeight))
	compoundImg.Draw(frame, 0, videoWidth)

	maxFps := pixelsPerSec
	err, downsampler := util.NewDownsampler(maxFps, fps)
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
