package assets

import (
	_ "embed"

	// "golang.org/x/image/font/sfnt"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed whitrabt.ttf
var fontBinary []byte

// func init() {
// 	parsedFont, err := truetype.Parse(fontBinary)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	face := truetype.NewFace(parsedFont, &truetype.Options{
// 		Size:    24, // font size in points
// 		DPI:     72, // dots per inch
// 		Hinting: font.HintingFull,
// 	})
// }

func WhitrabtFont(size float64) font.Face {
	f, err := truetype.Parse(fontBinary)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return face
}
