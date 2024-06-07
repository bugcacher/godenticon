package avatar

import (
	"os"
)

func byteSum(data []byte) uint8 {
	var sum uint8 = 0
	for _, b := range data {
		sum += b
	}
	return sum
}

func ensurePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		createErr := os.MkdirAll(path, 0755)
		if createErr != nil {
			return err
		}
	}
	return nil
}
