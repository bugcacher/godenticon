package avatar

type Algorithm int

const (
	ALGORITHM_1 Algorithm = iota
	ALGORITHM_2
)

type AvatarSize uint

const (
	AVATAR_SIZE_5 AvatarSize = 5
	AVATAR_SIZE_7 AvatarSize = 7
	AVATAR_SIZE_9 AvatarSize = 9
)

const (
	defaultFileName = "avatar.png"
)
