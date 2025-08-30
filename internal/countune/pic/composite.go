package pic

/**
 ** Purpose: A "countune composite" is a composition of multiple countune images,
 **          stitched together horizontally one after the other.
 **
 **/

import (
	// local
	"countube/internal/common"

	// standard
	"image"
)

type PicSupplier interface {
	Next() image.Image
}

func BuildCountuneStrip(picSupplier PicSupplier) *image.RGBA {

	var pics []image.Image

	for {
		pic := picSupplier.Next()
		if pic == nil {
			break
		}

		pics = append(pics, pic)
	}

	return common.StitchImagesHorizontally(pics)
}
