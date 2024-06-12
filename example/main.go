package main

import (
	"log"

	"github.com/bugcacher/godenticon/avatar"
)

func main() {
	avatar := avatar.New(
		"abhinavsingh",
		avatar.WithPixelPattern(avatar.PIXEL_PATTERN_5),
		avatar.WithAlgorithm(avatar.ALGORITHM_2),
		avatar.WithOutputDir("avatars"),
		avatar.WithDimension(200),
	)

	_, err := avatar.Generate()
	if err != nil {
		log.Default().Fatalf("failed to create avatar. Error: %v", err)
	}
}
