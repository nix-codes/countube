package common

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
)

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func EnsurePath(path string) {

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func FileExists(path string) bool {

	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func UrlExists(url string) bool {
	r, e := http.Head(url)
	return e == nil && r.StatusCode == 200
}

func DownloadFile(URL, fileName string) error {
	// get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	// create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the bytes to the fie
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func ConcatImages(leftImg image.Image, rightImg image.Image) *image.RGBA {

	lBounds := leftImg.Bounds()
	rBounds := rightImg.Bounds()
	maxHeight := Max(lBounds.Dy(), rBounds.Dy())

	outBounds := image.Rectangle{
		image.ZP,
		image.Point{lBounds.Dx() + rBounds.Dx(), maxHeight}}
	outRgba := image.NewRGBA(outBounds)
	ConcatImagesIntoTarget(leftImg, rightImg, outRgba)

	return outRgba
}

func ConcatImagesIntoTarget(leftImg image.Image, rightImg image.Image, targetImg *image.RGBA) {

	lBounds := leftImg.Bounds()
	rBounds := rightImg.Bounds()
	maxHeight := Max(lBounds.Dy(), rBounds.Dy())

	outYPos := (maxHeight - lBounds.Dy()) / 2
	outTmpBounds := image.Rectangle{
		image.Point{0, outYPos},
		image.Point{lBounds.Dx(), outYPos + lBounds.Dy()}}
	draw.Draw(targetImg, outTmpBounds, leftImg, image.ZP, draw.Src)

	outYPos = (maxHeight - rBounds.Dy()) / 2
	outTmpBounds = image.Rectangle{
		image.Point{lBounds.Dx(), outYPos},
		image.Point{lBounds.Dx() + rBounds.Dx(), outYPos + rBounds.Dy()}}

	draw.Draw(targetImg, outTmpBounds, rightImg, image.ZP, draw.Src)
}

func WritePngToFile(filePath string, img image.Image) {
	outFile, err := os.Create(filePath)
	CheckErr(err)
	defer outFile.Close()

	err = png.Encode(outFile, img)
	CheckErr(err)
}

func WriteJpegToFile(filePath string, img image.Image) {
	outFile, err := os.Create(filePath)
	CheckErr(err)
	defer outFile.Close()

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(outFile)

	jpeg.Encode(buf, img, &jpeg.Options{100})
	_, err = writer.Write(buf.Bytes())
	CheckErr(err)
}

func ReadImageFromFile(filePath string) image.Image {

	f, err := os.Open(filePath)
	CheckErr(err)
	defer f.Close()

	img, _, err := image.Decode(f)
	CheckErr(err)

	return img
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
