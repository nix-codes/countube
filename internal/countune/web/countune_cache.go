package web

import (
	// local
	"countube/internal/common"
	"countube/internal/countune"
	"countube/internal/countune/pic"

	// standard
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	// external
	"github.com/bits-and-blooms/bitset"
)

var countunePicLocalFilenameRegex, _ = regexp.Compile(CountunePicLocalFilenameRegexStr)

func DownloadMissingPicSpecs() {
	fmt.Println("Downloading Countune pic specs ...")
	var prevSpec *CountuneSpec

	for id := 0; id == 0 || prevSpec != nil; id++ {

		spec := specRepo.GetById(id)
		if spec == nil {
			spec = FetchPicSpec(id, prevSpec)
			if spec != nil {
				specRepo.Upsert(*spec)
			}
		}
		prevSpec = spec

		if id%10 == 0 {
			specRepo.Persist()
		}
	}
	specRepo.Persist()
	fmt.Println()
	fmt.Println("Finished downloading.")
}

func FetchPicSpec(id int, prevSpec *CountuneSpec) *CountuneSpec {
	fmt.Printf("Fetching specs for Countune pic#%d\r", id)
	time.Sleep(100 * time.Millisecond)
	img, _ := FetchPic(id)
	if img == nil {
		return nil
	}
	barCount, colors := pic.DetermineCountunePicSpec(img)

	startNum := 1
	if prevSpec != nil {
		startNum = (*prevSpec).StartNum + (*prevSpec).BarCount
	}

	return &CountuneSpec{
		Id:        id,
		StartNum:  startNum,
		BarCount:  barCount,
		BarColor1: colors[0],
		BarColor2: colors[1],
		BgColor:   colors[2],
	}
}

// This is deprecated. Instead of storing the images, we can just store the pic
// specs and generate them when needed
func SyncLocalPicCache() {
	fmt.Println("Updating local Countune cache...")
	common.EnsurePath(CountunePicCachePath)
	cachedPicNums, maxPicNum := checkPicsInLocalCache()
	downloadMissingPics(cachedPicNums, maxPicNum, CountunePicCachePath)
}

func ScanCountunePics(countunePicsPath string) []CountuneMeta {
	files, err := ioutil.ReadDir(countunePicsPath)
	if err != nil {
		log.Fatal(err)
	}

	picMetas := make([]CountuneMeta, len(files))

	for i, currentBar := 0, 1; ; i++ {
		filePath := GetPicFilePath(i)

		if !common.FileExists(filePath) {
			break
		}

		img := common.ReadImageFromFile(filePath)

		picMetas[i] = CountuneMeta{
			Id:     i,
			BarSeq: currentBar,
			Bars:   img.Bounds().Dx() / countune.CountunePicOriginalBarWidth,
		}

		currentBar += picMetas[i].Bars
	}

	return picMetas
}

func checkPicsInLocalCache() (*bitset.BitSet, int) {

	files, err := ioutil.ReadDir(CountunePicCachePath)
	if err != nil {
		log.Fatal(err)
	}
	cachedPicNums := bitset.New(10000)
	maxPicNum := -1

	if err == nil {

		for _, file := range files {
			picNum := parseCountunePicFileName(file.Name())

			if picNum >= 0 {
				cachedPicNums.Set(uint(picNum))

				if picNum > maxPicNum {
					maxPicNum = picNum
				}
			}
		}
	}

	return cachedPicNums, maxPicNum
}

func downloadMissingPics(cachedPicNums *bitset.BitSet, maxPicNum int, targetPath string) {

	// intermediate pics missing in the cache
	for i := 0; i < maxPicNum; i++ {
		if !cachedPicNums.Test(uint(i)) {
			downloadCountunePic(i, targetPath)
			time.Sleep(100 * time.Millisecond)
		}
	}

	// fetch pics after the last one found in the cache
	for i := maxPicNum + 1; common.UrlExists(getCountunePicUrl(i)); i++ {
		downloadCountunePic(i, targetPath)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("                                                        ")
	fmt.Println()
}

func downloadCountunePic(picNum int, targetPath string) {
	url := getCountunePicUrl(picNum)
	fileName := fmt.Sprintf("%s/%05d.png", targetPath, picNum)
	fmt.Printf("Downloading Countune pic#%d\r ...", picNum)
	err := common.DownloadFile(url, fileName)
	if err != nil {
		fmt.Println("\nFound a problem downloading file: ", err)
		fmt.Println("Aborting.")
		os.Exit(1)
	}
}

// func getCountunePicUrl(picNum int) string {
// 	return fmt.Sprintf(CountunePicUrlFmt, picNum)
// }

func parseCountunePicFileName(fileName string) int {

	picNumAsStr := countunePicLocalFilenameRegex.FindStringSubmatch(fileName)
	picNum, err := strconv.Atoi(picNumAsStr[1])
	if err != nil {
		return -1
	}

	return picNum
}

func GetPicFilePath(picNum int) string {
	fileName := fmt.Sprintf(CountunePicLocalFileNameFmt, picNum)
	return filepath.Join(CountunePicCachePath, fileName)
}

// func DetermineCountunePictureColors(img image.Image) [3]string {
// 	bounds := img.Bounds()
// 	width := bounds.Dx()
// 	barCount := width / countune.CountunePicOriginalBarWidth

// 	colors0 := getColorSamplesFromBar(img, 0, countune.CountunePicOriginalBarWidth)
// 	colors1 := getColorSamplesFromBar(img, 1, countune.CountunePicOriginalBarWidth)
// 	colors, err := determineCountuneColors(colors0[0], colors0[1], colors1[0], colors1[1])
// 	fmt.Println(colors0)
// 	fmt.Println(colors1)

// 	if err != nil {
// 		fmt.Println("checking from end bars")
// 		colors0 = getColorSamplesFromBar(img, barCount-1, countune.CountunePicOriginalBarWidth)
// 		colors1 = getColorSamplesFromBar(img, barCount-2, countune.CountunePicOriginalBarWidth)
// 		colors, err = determineCountuneColors(colors0[0], colors0[1], colors1[0], colors1[1])
// 	}

// 	if err != nil {
// 		fmt.Println("err!!!")
// 	}

// 	fmt.Println(colors)

// 	evenBarColor := colors[0]
// 	oddBarColor := colors[1]
// 	backgroundColor := colors[2]

// 	return [3]string{evenBarColor, oddBarColor, backgroundColor}
// }

// func determineCountuneColors(
// 	evenBarTopColor, evenBarBottomColor,
// 	oddBarTopColor, oddBarBottomColor string,
// ) ([]string, error) {

// 	// Count occurrences
// 	counts := make(map[string]int)
// 	colorsInput := []string{
// 		evenBarTopColor, evenBarBottomColor,
// 		oddBarTopColor, oddBarBottomColor,
// 	}
// 	for _, c := range colorsInput {
// 		counts[c]++
// 	}

// 	if len(counts) != 3 {
// 		return nil, fmt.Errorf("Cannot determine the colors effectively because we need 3 different colors in the samples.")
// 	}

// 	var result []string

// 	foundBackgroundColor := false
// 	for color, count := range counts {
// 		if count == 2 {
// 			result[2] = color
// 			foundBackgroundColor = true
// 			break
// 		}
// 	}

// 	if !foundBackgroundColor {
// 		return nil, fmt.Errorf("no color appears exactly 3 times")
// 	}

// 	if evenBarTopColor != result[2] {
// 		result[0] = evenBarTopColor
// 	} else {
// 		result[0] = evenBarBottomColor
// 	}

// 	if oddBarTopColor != result[2] {
// 		result[1] = oddBarTopColor
// 	} else {
// 		result[1] = oddBarBottomColor
// 	}

// 	return result, nil
// }

// func getColorSamplesFromBar(img image.Image, barIdx int, barWidth int) [2]string {
// 	height := img.Bounds().Dy()
// 	x := barIdx*barWidth + 1
// 	topColor := pixelColorHex(img, x, 1)
// 	bottomColor := pixelColorHex(img, x, height-1)

// 	return [2]string{topColor, bottomColor}
// }

// func pixelColorHex(img image.Image, x, y int) string {
// 	// we acocunt for the case where there's alpha channel
// 	c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
// 	return fmt.Sprintf("%02x%02x%02x", c.R, c.G, c.B)
// }
