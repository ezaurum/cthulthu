package paint

import (
	"github.com/nfnt/resize"
	"image"
)

func ResizeToInch(source image.Image, maxWidth float64, maxHeight float64, dpi int) image.Image {

	maxWidthPx := uint(InchToPixel(dpi, maxWidth))
	maxHeightPx := uint(InchToPixel(dpi, maxHeight))

	return ResizeWithRatio(source, maxWidthPx, maxHeightPx)
}

func CentimeterToPixel(dpi int, centi float64) int {
	return InchToPixel(dpi, CentimeterToInch(centi))
}

func InchToPixel(dpi int, maxWidth float64) int {
	return int(float64(dpi) * maxWidth)
}

func CentimeterToInch(centi float64) float64 {
	return centi * 0.393701
}

func Resize(source image.Image, width uint, height uint) image.Image {
	return resize.Resize(width, height, source, resize.Lanczos3)
}

func ResizeWithRatio(source image.Image, maxWidthPx uint, maxHeightPx uint) image.Image {
	width := source.Bounds().Dx()
	height := source.Bounds().Dy()
	if width > height {
		return resize.Resize(maxWidthPx, 0, source, resize.Lanczos3)
	} else {
		return resize.Resize(0, maxHeightPx, source, resize.Lanczos3)
	}

}
