package web

import (
	"countube/internal/common"
)

const (
	OutputPath                       = "./"
	CountunePicSpecStorageFile       = "./countune-specs.json"
	CountunePicCachePath             = OutputPath + "./countune-pics"
	CountunePicUrlFmt                = "https://www.countune.com/system/modules/xcountune/html/countunes/countune_%d.png"
	CountunePicLocalFilenameRegexStr = "([0-9]{5}).png"
	CountunePicLocalFileNameFmt      = "%05d.png"
)

type CountuneMeta struct {
	Id     int
	BarSeq int
	Bars   int
}

var specRepo *CountuneSpecRepo

func init() {
	var err error
	specRepo, err = NewCountuneSpecRepo()
	common.CheckErr(err)

}
