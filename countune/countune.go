package countune

import (
	"image/color"
)

const (
	CountunePicCachePath             = "./countune-pics"
	CountunePicUrlFmt                = "https://www.countune.com/system/modules/xcountune/html/countunes/countune_%d.png"
	CountunePicLocalFilenameRegexStr = "([0-9]{5}).png"
	CountunePicLocalFileNameFmt      = "%05d.png"

	CountunePicOriginalHeight        = 200
	CountunePicOriginalBarWidth      = 20
	CountunePicBarWidthToHeightRatio = int(CountunePicOriginalHeight / CountunePicOriginalBarWidth)
	CountunePicHeightToWidthRatio    = 3

	OutputPath                      = "./out"
	OutputRandomCountuneFilenameExt = ".countune.png"
	OutputTitlePicFilenameExt       = ".title.png"
	OutputFullVideoImageFilenameExt = ".scroll.png"
	OutputVideoFramesFilenameExt    = ".frames"

	FontPath = "resources/whitrabt.ttf"
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
