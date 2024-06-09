package avatar

type Algorithm int

const (
	ALGORITHM_1 Algorithm = iota
	ALGORITHM_2
)

type PatternSize uint

const (
	PATTERN_SIZE_5 PatternSize = 5
	PATTERN_SIZE_7 PatternSize = 7
	PATTERN_SIZE_9 PatternSize = 9
)

type Output int

const (
	OUTPUT_FILE Output = iota
	OUTPUT_BUFFER
)

const (
	defaultFileName = "avatar.png"
)
