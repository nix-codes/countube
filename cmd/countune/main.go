package main

import (
	"countube/internal/common"
	"countube/internal/countune/pic"
)

func main() {
	// spec := pic.NewCountuneCompositeSpec(1)
	// spec.AddPic(25, "000080", "bc8f8f", "eee8aa")
	// spec.AddPic(49, "ffff00", "adff2f", "ff0000")
	// spec.AddPic(36, "696969", "ff6347", "800080")
	// spec.AddPic(16, "ffa500", "0000ff", "000000")
	// spec.AddPic(25, "0000ff", "00ffff", "008080")

	// img := pic.GenerateCountuneFromSpec(spec)
	// common.WritePngToFile("test.png", img)

	// pic.Test()

	picSeqSpec := pic.PicSeqSpec{
		PicHeightInPixels:   200,
		PixelsPerUnit:       20,
		AmplitudeUnits:      10,
		StartNum:            1,
		StartHeightUnits:    5,
		StartDirection:      pic.UP,
		InitialBarPlacement: pic.ABOVE_WAVE,
	}

	picProvider := pic.NewCountunePicProvider(picSeqSpec)
	img := picProvider.NextRandomPic()
	common.WritePngToFile("test.png", img)
}
