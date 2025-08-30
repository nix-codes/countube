package generator

import (
	"fmt"
)

func BuildVideoGenCommand(vidCfg VideoConfig) string {
	imageFile := vidCfg.Name + OutputFullVideoImageFilenameExt
	videoFile := vidCfg.Name + ".mp4"

	cmd := "ffmpeg -loop 1 -i %s " +
		"-vf \"pad=iw+%d:ih:0:0:black, fps=%d, crop=%d:%d:x='if(lt(t,%d),0,(t-5)*%d)':y=0\" " +
		"-c:v libx264 -pix_fmt yuv420p -t %d %s"

	return fmt.Sprintf(cmd,
		imageFile,
		vidCfg.ScreenWidth,
		vidCfg.Fps,
		vidCfg.ScreenWidth,
		vidCfg.ScreenHeight,
		vidCfg.TitleDelay,
		vidCfg.BarWidth,
		vidCfg.VideoLen,
		videoFile)
}
