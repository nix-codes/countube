package main

import (
	// local
	"countube/internal/video/player"

	// standard
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	config := play.Config{
		ScreenWidth:  1920,
		ScreenHeight: 600,
		BarWidth:     60,
		ScrollSpeed:  2,
	}

	play.Play(config)
}
