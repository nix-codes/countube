package main

import (
	// local
	gen "countube/internal/video/generator"

	// standard
	"fmt"
	"image/color"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	vidCfg := gen.VideoConfig{
		Name:            "sample",
		Fps:             30,
		ScreenWidth:     1920,
		ScreenHeight:    1080,
		CountuneHeight:  800,
		BarWidth:        80,
		ScrollSpeed:     1,
		VideoLen:        30,
		BackgroundColor: color.Black,
		TitleUpperText: []string{
			"Band Name",
			"Some Album Name [1984]",
		},
		TitleLowerText: []string{
			"Video: Nix [Jun, 2022]",
			"Visual Art: Gerd Jansen [2009-]",
		},
		TitleDelay: 5,
		Texts: []gen.VideoText{
			{
				Text:         "1. First Track",
				StartSeconds: 0,
			},
		},
	}

	gen.PrepareImagesForVideo(vidCfg)
	ffmpegCmd := gen.BuildVideoGenCommand(vidCfg)

	fmt.Println()
	fmt.Println("Use the following command to generate the video:")
	fmt.Println(ffmpegCmd)
}
