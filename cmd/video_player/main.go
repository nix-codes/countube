package main

import (
	"countube/internal/countune/pic"

	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- Configuration ---
const (
	screenWidth  = 1080
	screenHeight = 200
	scrollSpeed  = 60 // pixels per second
)

// --- Image Stream Placeholder ---
var picSeqSpec = pic.PicSeqSpec{
	PicHeightInPixels:   screenHeight,
	PixelsPerUnit:       20,
	AmplitudeUnits:      10,
	StartNum:            1,
	StartHeightUnits:    5,
	StartDirection:      pic.UP,
	InitialBarPlacement: pic.ABOVE_WAVE,
}

var picProvider = pic.NewCountunePicProvider(picSeqSpec)

func nextImage() *ebiten.Image {
	// fmt.Println("nextImage()")
	img := picProvider.NextRandomPic()
	return ebiten.NewImageFromImage(img)
}

// --- Game State ---
type Game struct {
	imageQueue   []*ebiten.Image
	scrollOffset float64
	scrollSpeed  float64 // pixels per second
	lastUpdate   time.Time
}

func NewGame() *Game {
	g := &Game{
		imageQueue:   []*ebiten.Image{},
		scrollOffset: 0,
		scrollSpeed:  scrollSpeed,
	}

	// Pre-fill queue to cover screen width
	totalWidth := 0
	for totalWidth < screenWidth {
		img := nextImage()
		g.imageQueue = append(g.imageQueue, img)
		totalWidth += img.Bounds().Dx()
	}

	return g
}

func (g *Game) Update() error {
	now := time.Now()
	if g.lastUpdate.IsZero() {
		g.lastUpdate = now
		return nil
	}

	dt := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now
	g.scrollOffset += g.scrollSpeed * dt

	// Remove images that have fully scrolled past
	for len(g.imageQueue) > 0 && g.scrollOffset >= float64(g.imageQueue[0].Bounds().Dx()) {
		g.scrollOffset -= float64(g.imageQueue[0].Bounds().Dx())
		g.imageQueue = g.imageQueue[1:]
	}

	// Compute total width of images currently in the queue
	totalWidth := 0
	for _, img := range g.imageQueue {
		totalWidth += img.Bounds().Dx()
	}

	// Append new images if needed to fill the screen
	for float64(totalWidth)-g.scrollOffset < screenWidth {
		newImg := nextImage()
		g.imageQueue = append(g.imageQueue, newImg)
		totalWidth += newImg.Bounds().Dx()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	xPos := -g.scrollOffset
	for _, img := range g.imageQueue {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(xPos, 0)
		screen.DrawImage(img, op)
		xPos += float64(img.Bounds().Dx())
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Countune")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
