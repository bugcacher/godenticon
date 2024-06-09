package avatar

import (
	"image"
	"image/color"
	"math/rand"
)

type algoFunc func(img *image.RGBA, size int, colorToFill color.Color, darkMode bool)

var algoExecutorMap = map[Algorithm]algoFunc{
	ALGORITHM_1: algorithm_one,
	ALGORITHM_2: algorithm_two,
}

func algorithm_one(img *image.RGBA, size int, colorToFill color.Color, darkMode bool) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if y <= int(size)/2 {
				if rand.Float64() < 0.5 {
					img.Set(y, x, colorToFill)
				} else {
					img.Set(y, x, getBackgroundColor(darkMode))
				}
			} else {
				img.Set(y, x, img.At(int(size)-y-1, x))
			}
		}
	}
}

func algorithm_two(img *image.RGBA, size int, colorToFill color.Color, darkMode bool) {
	bounds := img.Bounds()
	for y := bounds.Max.Y; y >= 0; y-- {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x <= int(size)/2 {
				if rand.Float64() < 0.5 {
					img.Set(x, y, colorToFill)
				} else {
					img.Set(x, y, getBackgroundColor(darkMode))
				}
			} else {
				img.Set(x, y, img.At((int(size))-x-1, y))
			}
		}
	}
}

func getBackgroundColor(darkMode bool) color.Color {
	if darkMode {
		return color.Black
	}
	return color.White
}
