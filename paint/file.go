package paint

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

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
