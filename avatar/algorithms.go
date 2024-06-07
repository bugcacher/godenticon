package avatar

import (
	"image"
	"image/color"
	"math/rand"
)

type algoFunc func(image *image.RGBA, size int, colorToFill color.Color, darkMode bool)

var algoExecutorMap = map[Algorithm]algoFunc{
	ALGORITHM_1: algorithm_one,
	ALGORITHM_2: algorithm_two,
}

func algorithm_one(image *image.RGBA, size int, colorToFill color.Color, darkMode bool) {
	bounds := image.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if y <= int(size)/2 {
				if rand.Float64() < 0.5 {
					image.Set(y, x, colorToFill)
				} else {
					image.Set(y, x, getBackgroundColor(darkMode))
				}
			} else {
				image.Set(y, x, image.At(int(size)-y-1, x))
			}
		}
	}
}

func algorithm_two(image *image.RGBA, size int, colorToFill color.Color, darkMode bool) {
	bounds := image.Bounds()
	for y := bounds.Max.Y; y >= 0; y-- {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x <= int(size)/2 {
				if rand.Float64() < 0.5 {
					image.Set(x, y, colorToFill)
				} else {
					image.Set(x, y, getBackgroundColor(darkMode))
				}
			} else {
				image.Set(x, y, image.At((int(size))-x-1, y))
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
