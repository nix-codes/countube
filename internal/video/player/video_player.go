package play

import (
	"countube/internal/countune/pic"

	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	ScreenWidth  int
	ScreenHeight int
	BarWidth     int
	ScrollSpeed  float64 // bars per second
}

type video struct {
	config       Config
	imageQueue   []*ebiten.Image
	picSeq       *pic.OnTheFlyPicSeq
	scrollOffset float64
	scrollSpeed  float64 // pixels per second
	lastUpdate   time.Time
}

func Play(config Config) {
	video := newVideo(config)
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Countune")

	if err := ebiten.RunGame(video); err != nil {
		log.Fatal(err)
	}
}

func picSpecSupplier() *pic.CountunePicSpec {
	picSpec := pic.RandomPicSpec()
	return &picSpec
}

func (v *video) nextImage() *ebiten.Image {
	img := v.picSeq.Next()
	return ebiten.NewImageFromImage(img)
}

func newVideo(config Config) *video {
	var countuneSpec = pic.RandomCountuneSpec(config.ScreenHeight, config.BarWidth)
	var picSeq = pic.NewOnTheFlyPicSeq(countuneSpec, picSpecSupplier)

	v := &video{
		config:       config,
		imageQueue:   []*ebiten.Image{},
		scrollOffset: 0,
		scrollSpeed:  config.ScrollSpeed * float64(config.BarWidth),
		picSeq:       picSeq,
	}

	// Pre-fill queue to cover screen width
	totalWidth := 0
	for totalWidth < config.ScreenWidth {
		img := v.nextImage()
		v.imageQueue = append(v.imageQueue, img)
		totalWidth += img.Bounds().Dx()
	}

	return v
}

func (v *video) Update() error {
	now := time.Now()
	if v.lastUpdate.IsZero() {
		v.lastUpdate = now
		return nil
	}

	dt := now.Sub(v.lastUpdate).Seconds()
	v.lastUpdate = now
	v.scrollOffset += v.scrollSpeed * dt

	// Remove images that have fully scrolled past
	for len(v.imageQueue) > 0 && v.scrollOffset >= float64(v.imageQueue[0].Bounds().Dx()) {
		v.scrollOffset -= float64(v.imageQueue[0].Bounds().Dx())
		v.imageQueue = v.imageQueue[1:]
	}

	// Compute total width of images currently in the queue
	totalWidth := 0
	for _, img := range v.imageQueue {
		totalWidth += img.Bounds().Dx()
	}

	// Append new images if needed to fill the screen
	for float64(totalWidth)-v.scrollOffset < float64(v.config.ScreenWidth) {
		newImg := v.nextImage()
		v.imageQueue = append(v.imageQueue, newImg)
		totalWidth += newImg.Bounds().Dx()
	}

	return nil
}

func (v *video) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	xPos := -v.scrollOffset
	for _, img := range v.imageQueue {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(xPos, 0)
		screen.DrawImage(img, op)
		xPos += float64(img.Bounds().Dx())
	}
}

func (v *video) Layout(outsideWidth, outsideHeight int) (int, int) {
	return v.config.ScreenWidth, v.config.ScreenHeight
}
