package main

import (
	"log"

	"github.com/bugcacher/godenticon/avatar"
)

func main() {
	avatar := avatar.New(
		"abhinavsingh",
		avatar.WithSize(avatar.AVATAR_SIZE_5),
		avatar.WithAlgorithm(avatar.ALGORITHM_2),
		avatar.WithOutputDir("icons"),
		avatar.WithDarkMode(),
	)
	if err := avatar.GenerateAvatar(); err != nil {
		log.Default().Fatalf("failed to create avatar. Error: %v", err)
	}
}
