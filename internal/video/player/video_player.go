package play

import (
	// local
	"countube/internal/countune/pic"

	// standard
	"fmt"
	"image/color"
	"log"
	"time"

	// third-party
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	scrollSpeed  float64
	msg          string
	msgStart     time.Time
	lastUpdate   time.Time
	paused       bool
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
		scrollSpeed:  config.ScrollSpeed,
		picSeq:       picSeq,
		paused:       false,
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		v.paused = !v.paused
		v.lastUpdate = time.Now()
		v.msgStart = v.lastUpdate

		if v.paused {
			v.msg = "Paused"
		} else {
			v.msg = ""
		}

	}

	if v.paused {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		v.scrollSpeed += 0.5
		v.msg = fmt.Sprintf("Speed: %.01f bars/s", v.scrollSpeed)
		v.msgStart = time.Now()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		v.scrollSpeed -= 0.5
		if v.scrollSpeed < 0 {
			v.scrollSpeed = 0
		}
		v.msg = fmt.Sprintf("Speed: %.01f bars/s", v.scrollSpeed)
		v.msgStart = time.Now()
	}

	now := time.Now()
	if v.lastUpdate.IsZero() {
		v.lastUpdate = now
		return nil
	}

	dt := now.Sub(v.lastUpdate).Seconds()
	v.lastUpdate = now
	v.scrollOffset += v.scrollSpeed * float64(v.config.BarWidth) * dt

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

	if v.msg != "" && time.Since(v.msgStart) < 2*time.Second {
		ebitenutil.DebugPrintAt(screen, v.msg, 10, 10)
	}
}

func (v *video) Layout(outsideWidth, outsideHeight int) (int, int) {
	return v.config.ScreenWidth, v.config.ScreenHeight
}
