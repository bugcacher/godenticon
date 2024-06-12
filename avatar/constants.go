package avatar

type Algorithm int

const (
	ALGORITHM_1 Algorithm = iota
	ALGORITHM_2
)

type PixelPattern uint

const (
	PIXEL_PATTERN_5 PixelPattern = 5
	PIXEL_PATTERN_7 PixelPattern = 7
	PIXEL_PATTERN_9 PixelPattern = 9
)

type Output int

const (
	OUTPUT_FILE Output = iota
	OUTPUT_BUFFER
)

const (
	defaultFileName = "avatar.png"
)
