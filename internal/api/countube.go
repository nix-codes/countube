package api

import (
	// standard
	"image/color"
)

type VideoConfig struct {
	Name            string
	Fps             int
	CountuneHeight  int
	CountuneSpeed   float64
	BackgroundColor color.Color
	VideoWidth      int
	VideoHeight     int
	VideoLen        int
	Loop            bool
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
