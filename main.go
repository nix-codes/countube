package main

import (
	// local
	"countube/countune"

	// standard
	"fmt"
	"image/color"
)

func main() {

	vidCfg := countune.VideoConfig{
		Name:            "sample",
		Fps:             30,
		CountuneHeight:  800,
		CountuneSpeed:   1,
		VideoWidth:      1920,
		VideoHeight:     1080,
		VideoLen:        150,
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
		Texts: []countune.VideoText{
			{
				Text:         "1. First Track",
				StartSeconds: 0,
			},
			{
				Text:         "2. Second Track",
				StartSeconds: 45,
			},
			{
				Text:         "3. Third Track",
				StartSeconds: 120,
			},
		},
	}

	countune.UpdateLocalCache()
	countune.PrepareImagesForVideo(vidCfg)
	countune.GenerateImageScrollVideo(vidCfg)

	fmt.Println("done")
}
