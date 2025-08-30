package common

import (
	// standard
	"log"
	"os"
)

func EnsurePath(path string) {

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func FileExists(path string) bool {

	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
