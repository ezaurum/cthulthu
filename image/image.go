package image

import (
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/jpeg"
	"log"
	"os"
	"image/png"
)

func MakeQR(url string, size int) image.Image {
	qrCode, err := qrcode.New(url, qrcode.Low)
	if nil == err {
		return qrCode.Image(size)
	}
	return nil
}

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

func CreateJPEG(outputFilename string, canvas image.Image) {
	third, err := os.Create(outputFilename)
	defer third.Close()
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}
	jpeg.Encode(third, canvas, &jpeg.Options{jpeg.DefaultQuality})
}

func LoadPNG(filePath string) image.Image {

	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal("open:", err)
	}
	img, err := png.Decode(imgFile)
	if err != nil {
		log.Fatal("decode:", err)
	}
	if nil == img {
		log.Fatal("image: image is nul")
	}

	return img
}

func LoadJPEG(filePath string) image.Image {

	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Println("open:", err)
	}
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		log.Println("decode:", err)
	}
	if nil == img {
		log.Println("image: image is nul")
	}

	return img
}
