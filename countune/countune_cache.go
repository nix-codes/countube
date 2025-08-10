package countune

import (
	// local
	"countube/common"

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

func UpdateLocalCache() {
	fmt.Println("Updating local Countune cache...")
	common.EnsurePath(CountunePicCachePath)
	cachedPicNums, maxPicNum := checkPicsInLocalCache()
	downloadMissingPics(cachedPicNums, maxPicNum, CountunePicCachePath)
}

func scanCountunePics(countunePicsPath string) []countuneMeta {
	files, err := ioutil.ReadDir(countunePicsPath)
	if err != nil {
		log.Fatal(err)
	}

	picMetas := make([]countuneMeta, len(files))

	for i, currentBar := 0, 1; ; i++ {
		filePath := getPicFilePath(i)

		if !common.FileExists(filePath) {
			break
		}

		img := common.ReadImageFromFile(filePath)

		picMetas[i] = countuneMeta{
			id:     i,
			barSeq: currentBar,
			bars:   img.Bounds().Dx() / CountunePicOriginalBarWidth,
		}

		currentBar += picMetas[i].bars
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

func getCountunePicUrl(picNum int) string {
	return fmt.Sprintf(CountunePicUrlFmt, picNum)
}

func parseCountunePicFileName(fileName string) int {

	picNumAsStr := countunePicLocalFilenameRegex.FindStringSubmatch(fileName)
	picNum, err := strconv.Atoi(picNumAsStr[1])
	if err != nil {
		return -1
	}

	return picNum
}

func getPicFilePath(picNum int) string {
	fileName := fmt.Sprintf(CountunePicLocalFileNameFmt, picNum)
	return filepath.Join(CountunePicCachePath, fileName)
}
