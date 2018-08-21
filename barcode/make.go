package barcode

import (
	"image"
	"image/draw"
	"github.com/ezaurum/cthulthu/paint"
	"github.com/skip2/go-qrcode"
	"errors"
	"fmt"
	"image/color"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode"
)

const (
	ImageHeight      = 120
	ImagePadding     = 10
	DefaultMMSWidth  = 320
	DefaultMMSHeight = 1138
)

//640px by 1138px
func MakeMMSBarCodeFile(codeString string, fileName string, defaultImage image.Image) (error, bool) {
	cs, e := code128.Encode(codeString)
	if nil != e {
		return e, false
	}

	var backgroundBounds image.Rectangle
	var canvas draw.Image
	if nil != defaultImage {
		backgroundBounds = defaultImage.Bounds()
		i := cs.Bounds().Dx() + ImagePadding
		if backgroundBounds.Dx() < i {
			return errors.New(fmt.Sprintf("default images is too small, need bigger then %v", i)), false
		}
		canvas = image.NewRGBA(backgroundBounds)
	} else {
		backgroundBounds = image.Rect(0, 0, DefaultMMSWidth, ImageHeight)
		canvas = image.NewRGBA(backgroundBounds)
	}

	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.White}, image.ZP, draw.Src)

	barCode, err := barcode.Scale(cs, DefaultMMSWidth, ImageHeight)
	if nil != err {
		return err, false
	}

	barcodeWidth := barCode.Bounds().Dx()
	paddingA := (backgroundBounds.Dx() - int(barcodeWidth)) / 2

	maxPoint := backgroundBounds.Max
	barCodeBounds := image.Rect(paddingA,
		maxPoint.Y-ImageHeight, paddingA+barcodeWidth, maxPoint.Y)

	if nil != defaultImage {
		draw.Draw(canvas, backgroundBounds, defaultImage, image.ZP, draw.Src)
	}
	draw.Draw(canvas, barCodeBounds, barCode, image.ZP, draw.Src)

	paint.CreateJPEG(fileName, canvas)

	return nil, true
}

func MakeBarCodeFile(codeString string, fileName string) (error, bool) {
	img, err := MakeBarCode(codeString)
	if nil != err {
		return err, false
	}

	paint.CreateJPEG(fileName, img)

	return nil, true
}

func MakeBarCode(codeString string) (image.Image, error) {
	cs, e := code128.Encode(codeString)
	if nil != e {
		return nil, e
	}
	backgroundBounds := image.Rect(0, 0, 220, 12)
	canvas := image.NewRGBA(backgroundBounds)
	barCode := paint.Resize(cs, 200, 10)
	barCodeBounds := image.Rect(10, 1, 210, 11)
	draw.Draw(canvas, barCodeBounds, barCode, image.ZP, draw.Src)

	return canvas, nil
}

func MakeQR(url string, size int) image.Image {
	qrCode, err := qrcode.New(url, qrcode.Low)
	if nil == err {
		return qrCode.Image(size)
	}
	return nil
}
