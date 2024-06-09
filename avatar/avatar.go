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
	value       string
	path        string
	patternSize PatternSize
	algo        Algorithm
	darkMode    bool
	outputType  Output
	dimension   uint
	image       *image.RGBA
}

type AvatarResult struct {
	// Path contains the filepath of the avatar generated.
	// Path will be empty if Output type is OUTPUT_BUFFER
	Path string
	// Buffer contains the generate avatar buffer.
	// Buffer will be nil if Output type in OUTPUT_FILE
	Buffer *bytes.Buffer
}

// New returns an Avatar object which can be used to generate an identicon.
func New(value string, opts ...CreateOption) *Avatar {
	avatar := &Avatar{
		value:       value,
		patternSize: PATTERN_SIZE_5,
		algo:        ALGORITHM_1,
		outputType:  OUTPUT_FILE,
		dimension:   100,
	}
	for _, opt := range opts {
		opt(avatar)
	}
	return avatar
}

// WithPatternSize sets the PatternSize of the generated avatar.
// Pattern size is used to set the base image pixel pattern of the avatar.
// PATTERN_SIZE_5 creates an avatar of 5 * 5 pixel pattern.
// PatternSize is different from Dimension and is only used to set the base image pattern size.
func WithPatternSize(size PatternSize) func(a *Avatar) {
	return func(a *Avatar) {
		a.patternSize = size
	}
}

// WithOutputDir sets the directory path of the generate avatar image file.
// WithOutputDir will not have any effect if Output type is OUTPUT_BUFFER
func WithOutputDir(path string) func(a *Avatar) {
	if err := ensurePath(path); err != nil {
		log.Default().Fatal("Invalid path given")
	}
	return func(a *Avatar) {
		a.path = path
	}
}

// WithAlgorithm sets the algorithm used for generating the avatar
func WithAlgorithm(algo Algorithm) func(a *Avatar) {
	return func(a *Avatar) {
		a.algo = algo
	}
}

// WithDarkMode is used to generate the avatar in dark mode.
// In dark mode avatar background is of Black color instead of White color (default).
func WithDarkMode() func(a *Avatar) {
	return func(a *Avatar) {
		a.darkMode = true
	}
}

// WithOutputType sets the Output type. Avatar can be saved in
// a file or in buffer.
func WithOutputType(outputType Output) func(a *Avatar) {
	return func(a *Avatar) {
		a.outputType = outputType
	}
}

// WithDimension sets the dimensions (height * width) of the generated Avatar.
func WithDimension(dimension uint) func(a *Avatar) {
	return func(a *Avatar) {
		a.dimension = dimension
	}
}

// GenerateAvatar generates an unique avatar for the given value.
func (av *Avatar) GenerateAvatar() (*AvatarResult, error) {
	hash := sha256.Sum256([]byte(av.value))
	seed := binary.BigEndian.Uint32(hash[:])
	rand.Seed(int64(seed))

	r := uint8(uint64(byteSum(hash[0:8])) % 256)
	g := uint8(uint64(byteSum(hash[8:16])) % 256)
	b := uint8(uint64(byteSum(hash[16:24])) % 256)
	a := uint8(uint64(byteSum(hash[24:32])) % 256)
	avatarColor := color.RGBA{r, g, b, a}

	height, width := av.patternSize, av.patternSize
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
		return &AvatarResult{Path: filePath}, nil
	case OUTPUT_BUFFER:
		return &AvatarResult{Buffer: &buf}, nil
	}

	return nil, ErrUnknownOutputType

}

func (av *Avatar) applyAlgorithm(colorToFill color.Color, darkMode bool) {
	algoFunc := algoExecutorMap[av.algo]
	algoFunc(av.image, int(av.patternSize), colorToFill, darkMode)
}

func (av *Avatar) scaleImage() {
	scaledImage := image.NewRGBA(image.Rect(0, 0, int(av.dimension), int(av.dimension)))
	draw.NearestNeighbor.Scale(scaledImage, scaledImage.Bounds(), av.image, av.image.Bounds(), draw.Over, nil)
	av.image = scaledImage
}

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
