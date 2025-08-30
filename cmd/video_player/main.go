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
		ScreenHeight: 400,
		BarWidth:     40,
		ScrollSpeed:  1,
	}

	play.Play(config)
}
