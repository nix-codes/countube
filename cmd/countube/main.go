package main

import (
	// local
	"countube/internal/api"
	"countube/internal/video"

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

	vidCfg := api.VideoConfig{
		Name:            "sample",
		Fps:             30,
		CountuneHeight:  800,
		CountuneSpeed:   1,
		VideoWidth:      1920,
		VideoHeight:     1080,
		VideoLen:        90,
		BackgroundColor: color.Black,
		Loop:            true, // making it true will ignore the title-related params
		TitleUpperText: []string{
			"Band Name",
			"Some Album Name [1984]",
		},
		TitleLowerText: []string{
			"Video: Nix [Jun, 2022]",
			"Visual Art: Gerd Jansen [2009-]",
		},
		TitleDelay: 5,
		Texts: []api.VideoText{
			{
				Text:         "1. First Track",
				StartSeconds: 0,
			},
		},
	}

	//countune.UpdateLocalCache()
	video.PrepareImagesForVideo(vidCfg)
	// countune.GenerateImageScrollVideo(vidCfg)

	fmt.Println("done")
}
