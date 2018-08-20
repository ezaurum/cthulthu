package barcode

import (
	"image"
	"github.com/boombuler/barcode/code128"
	"image/draw"
	"github.com/ezaurum/cthulthu/paint"
	"github.com/skip2/go-qrcode"
	"errors"
	"fmt"
	"image/color"
	"github.com/nfnt/resize"
)

const (
	ImageHeight      = 200
	ImagePadding     = 20
	DefaultMMSWidth  = 640
	DefaultMMSHeight = 1138
)

//640px by 1138px
func MakeMMSBarCodeFile(codeString string, fileName string, defaultImage image.Image) (error, bool) {
	cs, e := code128.Encode(codeString)
	if nil != e {
		return e, false
	}
	//paint.CreateJPEG(cs.Content()+".jpg", paint.ResizeA(cs, uint(cs.Bounds().Dx()*2), ImageHeight, resize.NearestNeighbor))

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

	doubleSize := uint(cs.Bounds().Dx() * 4)
	barCode := paint.ResizeA(cs, doubleSize, ImageHeight, resize.NearestNeighbor)

	paddingA := (backgroundBounds.Dx() - int(doubleSize)) / 2
	paddingB := (backgroundBounds.Dx() - int(doubleSize)) - paddingA

	minPoint := backgroundBounds.Min
	maxPoint := backgroundBounds.Max
	barCodeBounds := image.Rect(minPoint.X+paddingA,
		maxPoint.Y-ImageHeight, maxPoint.X-paddingB, maxPoint.Y)

	paint.CreateJPEG("t0"+fileName,
		paint.ResizeA(cs, doubleSize, uint(barCodeBounds.Dy()), resize.Bicubic))

	paint.CreateJPEG("t1"+fileName,
		paint.ResizeA(cs, doubleSize, uint(barCodeBounds.Dy()), resize.Bilinear))

	paint.CreateJPEG("t2"+fileName,
		paint.ResizeA(cs, doubleSize, uint(barCodeBounds.Dy()), resize.NearestNeighbor))

	paint.CreateJPEG("t3"+fileName,
		paint.ResizeA(cs, doubleSize, uint(barCodeBounds.Dy()), resize.Lanczos3))

	if nil != defaultImage {
		draw.Draw(canvas, backgroundBounds, defaultImage, image.ZP, draw.Src)
	}
	draw.Draw(canvas, barCodeBounds, barCode, image.ZP, draw.Src)

	paint.CreateJPEG(cs.Content()+".jpg", canvas)
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
