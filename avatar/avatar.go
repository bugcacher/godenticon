// Package avatar provides functionality to create GitHub-like avatars (identicons).
// It allows customization of the avatar's pattern size, algorithm, output type, dimension, and color mode.
package avatar

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

type CreateOption func(a *Avatar)

type Avatar struct {
	value        string
	path         string
	dimension    uint
	darkMode     bool
	pixelPattern PixelPattern
	algo         Algorithm
	outputType   Output
	image        *image.RGBA
}

// AvatarResult contains the result of an avatar generation process.
type AvatarResult struct {
	// FilePath contains the file path of the generated avatar image.
	// FilePath will be empty if the OutputType is OutputBuffer.
	FilePath string
	// Buffer contains the generated avatar image as a byte buffer.
	// Buffer will be nil if the OutputType is OutputFile.
	Buffer *bytes.Buffer
}

// New creates and returns a new Avatar object with the specified value and options.
func New(value string, opts ...CreateOption) *Avatar {
	avatar := &Avatar{
		value:        value,
		pixelPattern: PIXEL_PATTERN_5,
		algo:         ALGORITHM_1,
		outputType:   OUTPUT_FILE,
		dimension:    100,
	}
	for _, opt := range opts {
		opt(avatar)
	}
	return avatar
}

// WithPixelPattern sets the pixel pattern size of the generated avatar.
// Pixel pattern size defines the base image pixel pattern of the avatar.
// For example, PIXEL_PATTERN_5 creates an avatar with a 5x5 pixel pattern.
// PixelPattern is different from Dimension and is only used to set the base pixel pattern size.
func WithPixelPattern(pixelPattern PixelPattern) func(a *Avatar) {
	return func(a *Avatar) {
		a.pixelPattern = pixelPattern
	}
}

// WithOutputPath sets the directory path for the generated avatar image file.
// This option is ignored if the output type is OutputBuffer.
func WithOutputDir(path string) func(a *Avatar) {
	if err := ensurePath(path); err != nil {
		log.Default().Fatal("Invalid path given")
	}
	return func(a *Avatar) {
		a.path = path
	}
}

// WithAlgorithm sets the algorithm used for generating the avatar.
func WithAlgorithm(algo Algorithm) func(a *Avatar) {
	return func(a *Avatar) {
		a.algo = algo
	}
}

// WithDarkMode enables dark mode for the avatar, setting the background color to black.
func WithDarkMode() func(a *Avatar) {
	return func(a *Avatar) {
		a.darkMode = true
	}
}

// WithOutputType sets the output type for the generated avatar.
// The avatar can be saved to a file or stored in a buffer.
func WithOutputType(outputType Output) func(a *Avatar) {
	return func(a *Avatar) {
		a.outputType = outputType
	}
}

// WithDimension sets the dimensions (height and width) of the generated avatar.
func WithDimension(dimension uint) func(a *Avatar) {
	return func(a *Avatar) {
		a.dimension = dimension
	}
}

// Generate creates a unique avatar for the given value based on the Avatar configuration.
func (av *Avatar) Generate() (*AvatarResult, error) {
	hash := sha256.Sum256([]byte(av.value))
	seed := binary.BigEndian.Uint32(hash[:])
	rand.Seed(int64(seed))

	r := uint8(uint64(byteSum(hash[0:8])) % 256)
	g := uint8(uint64(byteSum(hash[8:16])) % 256)
	b := uint8(uint64(byteSum(hash[16:24])) % 256)
	a := uint8(uint64(byteSum(hash[24:32])) % 256)
	avatarColor := color.RGBA{r, g, b, a}

	height, width := av.pixelPattern, av.pixelPattern
	av.image = image.NewRGBA(image.Rect(0, 0, int(height), int(width)))

	av.applyAlgorithm(avatarColor, av.darkMode)

	av.scaleImage()

	var buf bytes.Buffer
	if err := png.Encode(&buf, av.image); err != nil {
		return nil, err
	}

	switch av.outputType {
	case OUTPUT_FILE:
		filePath, err := av.saveToFile()
		if err != nil {
			return nil, err
		}
		return &AvatarResult{FilePath: filePath}, nil
	case OUTPUT_BUFFER:
		return &AvatarResult{Buffer: &buf}, nil
	}

	return nil, ErrUnknownOutputType
}

// applyAlgorithm applies the selected algorithm to generate the avatar's pixel pattern.
func (av *Avatar) applyAlgorithm(colorToFill color.Color, darkMode bool) {
	algoFunc := algoExecutorMap[av.algo]
	algoFunc(av.image, int(av.pixelPattern), colorToFill, darkMode)
}

// scaleImage scales the base image to the desired dimensions.
func (av *Avatar) scaleImage() {
	scaledImage := image.NewRGBA(image.Rect(0, 0, int(av.dimension), int(av.dimension)))
	draw.NearestNeighbor.Scale(scaledImage, scaledImage.Bounds(), av.image, av.image.Bounds(), draw.Over, nil)
	av.image = scaledImage
}

// saveToFile saves the generated avatar image to a file and returns the file path.
func (av *Avatar) saveToFile() (string, error) {
	outputPath := filepath.Join(av.path, defaultFileName)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	if err := png.Encode(outFile, av.image); err != nil {
		return "", err
	}
	return outputPath, nil
}
