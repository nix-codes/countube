package common

import (
	// standard
	"log"
)

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
