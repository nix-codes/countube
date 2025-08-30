package generator

import (
	"image/color"
)

type VideoConfig struct {
	Name            string
	ScreenWidth     int
	ScreenHeight    int
	BackgroundColor color.Color
	CountuneHeight  int
	BarWidth        int
	ScrollSpeed     float64 // bars per second
	Fps             int
	VideoLen        int
	TitleUpperText  []string
	TitleLowerText  []string
	TitleDelay      int
	Texts           []VideoText
}

type VideoText struct {
	Text         string
	StartSeconds int
}

const (
	OutputPath                      = "./out"
	OutputRandomCountuneFilenameExt = ".countune.png"
	OutputTitlePicFilenameExt       = ".title.png"
	OutputFullVideoImageFilenameExt = ".scroll.png"
	OutputVideoFramesFilenameExt    = ".frames"
)
