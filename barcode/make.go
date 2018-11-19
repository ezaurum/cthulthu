package barcode

import (
	"image"
	"image/draw"
	"github.com/ezaurum/cthulthu/paint"
	"github.com/skip2/go-qrcode"
	"errors"
	"fmt"
	"github.com/boombuler/barcode/twooffive"
	"github.com/boombuler/barcode"
	"github.com/golang/freetype"
	"io/ioutil"
	"log"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	ImageHeight      = 120
	DefaultMMSWidth  = 320
)

var (
	FontName = "luxisr.ttf"
	Font     *truetype.Font
)

// initialize font
func InitializeFont() {

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(FontName)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	Font = f
}

//640px by 1138px
func MakeMMSBarCodeFile(codeString string,
	fileName string, defaultImage image.Image,
	withCodeString bool) (error, bool) {

	var barCode image.Image
	var err error
	if withCodeString {
		barCode, err = MakeBarCodeWithString(codeString)
	} else {
		barCode, err = MakeBarCode(codeString)
	}
	if nil != err {
		return err, false
	}

	var backgroundBounds image.Rectangle
	var canvas draw.Image
	if nil != defaultImage {
		backgroundBounds = defaultImage.Bounds()
		if backgroundBounds.Dx() < DefaultMMSWidth {
			return errors.New(fmt.Sprintf("default images is too small, need bigger then %v", DefaultMMSWidth)), false
		}
		canvas = image.NewRGBA(backgroundBounds)
	} else {
		backgroundBounds = image.Rect(0, 0, DefaultMMSWidth, ImageHeight)
		canvas = image.NewRGBA(backgroundBounds)
		draw.Draw(canvas, canvas.Bounds(), image.White, image.ZP, draw.Src)
	}

	barcodeWidth := barCode.Bounds().Dx()
	paddingA := (backgroundBounds.Dx() - int(barcodeWidth)) / 2

	maxPoint := backgroundBounds.Max
	var barCodeBounds image.Rectangle
	if nil != defaultImage {
		barCodeBounds = image.Rect(paddingA,
			maxPoint.Y-ImageHeight, paddingA+barcodeWidth, maxPoint.Y)
	} else {
		barCodeBounds = image.Rect(paddingA,
			maxPoint.Y-ImageHeight, paddingA+barcodeWidth, maxPoint.Y)
	}

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
	cs, e := twooffive.Encode(codeString, false)
	if nil != e {
		return nil, e
	}

	barCode, err := barcode.Scale(cs, DefaultMMSWidth -20, ImageHeight)
	if nil != err {
		return nil, err
	}

	return barCode, nil
}

func MakeBarCodeWithString(codeString string) (image.Image, error) {

	code, e := MakeBarCode(codeString)

	if nil != e {
		return code, e
	}

	paddingLeft := 10
	bounds := code.Bounds()
	clip := image.Rect(bounds.Min.X+paddingLeft, bounds.Max.Y-50, bounds.Max.X-paddingLeft, bounds.Max.Y)
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, code, image.ZP, draw.Src)
	draw.Draw(canvas, clip, image.White, image.ZP, draw.Src)

	barCodeBounds := image.Rect(bounds.Min.X,
		bounds.Min.Y, bounds.Max.X, bounds.Max.Y-40)

	draw.Draw(canvas, barCodeBounds, code, image.ZP, draw.Src)

	fontSize := 24.0
	dpi := 72.0

	c := freetype.NewContext()
	c.SetDPI(dpi)
	if nil == Font {
		InitializeFont()
	}
	c.SetFont(Font)
	c.SetFontSize(fontSize)
	c.SetClip(clip)
	c.SetDst(canvas)
	c.SetSrc(image.Black)
	hinting := "none"
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	paddingTop := 10

	pt := freetype.Pt(paddingLeft*2, clip.Min.Y+int(c.PointToFixed(fontSize)>>6)+paddingTop)
	_, err := c.DrawString(codeString, pt)
	if err != nil {
		return canvas, err
	}

	return canvas, e
}

func MakeQR(url string, size int) image.Image {
	qrCode, err := qrcode.New(url, qrcode.Low)
	if nil == err {
		return qrCode.Image(size)
	}
	return nil
}
