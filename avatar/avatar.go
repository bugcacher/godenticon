package avatar

import (
	"crypto/sha256"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

type CreateOption func(a *Avatar)

type Avatar struct {
	value    string
	path     string
	size     AvatarSize
	algo     Algorithm
	image    *image.RGBA
	darkMode bool
}

func New(value string, opts ...CreateOption) *Avatar {
	avatar := &Avatar{
		value: value,
		size:  AVATAR_SIZE_5,
		algo:  ALGORITHM_1,
	}
	for _, opt := range opts {
		opt(avatar)
	}
	return avatar
}

func WithSize(size AvatarSize) func(a *Avatar) {
	return func(a *Avatar) {
		a.size = size
	}
}

func WithOutputDir(path string) func(a *Avatar) {
	if err := ensurePath(path); err != nil {
		log.Default().Fatal("Invalid path given")
	}
	return func(a *Avatar) {
		a.path = path
	}
}

func WithAlgorithm(algo Algorithm) func(a *Avatar) {
	return func(a *Avatar) {
		a.algo = algo
	}
}

func WithDarkMode() func(a *Avatar) {
	return func(a *Avatar) {
		a.darkMode = true
	}
}

func (av *Avatar) GenerateAvatar() error {
	hash := sha256.Sum256([]byte(av.value))
	seed := binary.BigEndian.Uint32(hash[:])
	rand.Seed(int64(seed))

	r := uint8(uint64(byteSum(hash[0:8])) % 256)
	g := uint8(uint64(byteSum(hash[8:16])) % 256)
	b := uint8(uint64(byteSum(hash[16:24])) % 256)
	a := uint8(uint64(byteSum(hash[24:32])) % 256)
	avatarColor := color.RGBA{r, g, b, a}

	height, width := av.size, av.size
	av.image = image.NewRGBA(image.Rect(0, 0, int(height), int(width)))

	av.applyAlgorithm(avatarColor, av.darkMode)

	return av.saveToFile()
}

func (av *Avatar) applyAlgorithm(colorToFill color.Color, darkMode bool) {
	algoFunc := algoExecutorMap[av.algo]
	algoFunc(av.image, int(av.size), colorToFill, darkMode)
}

func (av *Avatar) saveToFile() error {
	outputPath := filepath.Join(av.path, defaultFileName)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	if err := png.Encode(outFile, av.image); err != nil {
		return err
	}
	return nil
}
