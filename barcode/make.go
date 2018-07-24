package barcode

import (
	"image"
	"github.com/boombuler/barcode/code128"
	"image/draw"
	"github.com/ezaurum/cthulthu/paint"
	"github.com/skip2/go-qrcode"
)

func MakeMMSBarCodeFile(codeString string, fileName string, defaultImage image.Image) (error, bool) {
	cs, e := code128.Encode(codeString)
	if nil != e {
		return e, false
	}

	backgroundBounds := image.Rect(0, 0, 320, 480)
	canvas := image.NewRGBA(backgroundBounds)

	if nil != defaultImage {
		defaultImage = paint.Resize(defaultImage, 320, 480)
	}

	barCode := paint.Resize(cs, 280, 100)

	barCodeBounds := image.Rect(20, 370, 300, 480)
	if nil != defaultImage {
		draw.Draw(canvas, backgroundBounds, defaultImage, image.ZP, draw.Src)
	}
	draw.Draw(canvas, barCodeBounds, barCode, image.ZP, draw.Src)

	paint.CreateJPEG(fileName, canvas)

	return nil, true
}

func MakeBarCodeFile(codeString string, fileName string) (error, bool) {
	cs, e := code128.Encode(codeString)
	if nil != e {
		return e, false
	}

	backgroundBounds := image.Rect(0, 0, 220, 12)
	canvas := image.NewRGBA(backgroundBounds)

	barCode := paint.Resize(cs, 200, 10)

	barCodeBounds := image.Rect(10, 1, 210, 11)
	draw.Draw(canvas, barCodeBounds, barCode, image.ZP, draw.Src)

	paint.CreateJPEG(fileName, canvas)

	return nil, true
}

func MakeQR(url string, size int) image.Image {
	qrCode, err := qrcode.New(url, qrcode.Low)
	if nil == err {
		return qrCode.Image(size)
	}
	return nil
}

