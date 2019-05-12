package barcode

import (
	"errors"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/twooffive"
	"github.com/ezaurum/cthulthu/paint"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/font"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
)

const (
	ImageHeight     = 240
	DefaultMMSWidth = 640
)

var (
	FontName = "luxisr.ttf"
	Font     *truetype.Font
)

// initialize font
func InitializeFont() error {

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(FontName)
	if err != nil {
		log.Println(err)
		return err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return err
	}

	Font = f
	return nil
}

//640px by 1138px
func MakeMMSBarCodeFile(codeString string,
	fileName string, defaultImage image.Image,
	settings map[string]interface{}) (error, bool) {

	// set default value
	withCodeString := false
	barCodeType := "128"
	paddingLeft := 130
	paddingBottom := 10
	fontSize := 32.0
	dpi := 72.0

	// override settings
	if nil != settings {
		if i, b := settings["withCodeString"]; b {
			withCodeString = i.(bool)
		}
		if i, b := settings["barcodeType"]; b {
			barCodeType = i.(string)
		}
		if i, b := settings["paddingLeft"]; b {
			paddingLeft = i.(int)
		}
		if i, b := settings["paddingBottom"]; b {
			paddingBottom = i.(int)
		}
		if i, b := settings["fontSize"]; b {
			fontSize = i.(float64)
		}
		if i, b := settings["dpi"]; b {
			dpi = i.(float64)
		}
	}

	var barCode image.Image
	var err error
	if withCodeString {
		barCode, err = MakeBarCodeWithString(codeString, barCodeType, paddingLeft, paddingBottom, fontSize, dpi)
	} else {
		barCode, err = MakeBarCode(codeString, barCodeType)
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
			maxPoint.Y-ImageHeight + 30, paddingA+barcodeWidth, maxPoint.Y)
	} else {
		barCodeBounds = image.Rect(paddingA,
			maxPoint.Y-ImageHeight + 30, paddingA+barcodeWidth, maxPoint.Y)
	}

	if nil != defaultImage {
		draw.Draw(canvas, backgroundBounds, defaultImage, image.ZP, draw.Src)
	}
	draw.Draw(canvas, barCodeBounds, barCode, image.ZP, draw.Src)

	paint.CreateJPEG(fileName, canvas)

	return nil, true
}

func MakeBarCodeFile(codeString string, fileName string, barcodeType string) (error, bool) {
	img, err := MakeBarCode(codeString, barcodeType)
	if nil != err {
		return err, false
	}

	paint.CreateJPEG(fileName, img)

	return nil, true
}

func MakeBarCode(codeString string, barcodeType string) (image.Image, error) {
	var cs barcode.Barcode
	var e error
	switch barcodeType {
	case "39":
		cs, e = code39.Encode(codeString, false, false)
	case "128":
		cs, e = code128.Encode(codeString)
	case "25":
		cs, e = twooffive.Encode(codeString, false)
	}
	if nil != e {
		return nil, e
	}

	barCode, err := barcode.Scale(cs, DefaultMMSWidth-40, ImageHeight-60)
	if nil != err {
		return nil, err
	}

	return barCode, nil
}

func MakeBarCodeWithString(codeString string, barcodeType string, paddingLeft int, paddingBottom int, fontSize float64, dpi float64) (image.Image, error) {

	code, e := MakeBarCode(codeString, barcodeType)

	if nil != e {
		return code, e
	}

	if nil == Font {
		fmt.Println("font is nil")
		e := InitializeFont()
		if nil != e {
			fmt.Println("font load error")
			return nil, e
		}
	}

	bounds := code.Bounds()
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, code, image.ZP, draw.Src)

	barCodeBounds := image.Rect(bounds.Min.X,
		bounds.Min.Y, bounds.Max.X, bounds.Max.Y)

	draw.Draw(canvas, barCodeBounds, code, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(Font)
	c.SetFontSize(fontSize)
	c.SetDst(canvas)
	c.SetSrc(image.Black)
	hinting := "none"
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	height := int(c.PointToFixed(fontSize) >> 6)
	clip := image.Rect(bounds.Min.X+paddingLeft, bounds.Max.Y-paddingBottom-height, bounds.Max.X-paddingLeft, bounds.Max.Y)
	c.SetClip(clip)

	draw.Draw(canvas, clip, image.White, image.ZP, draw.Src)

	pt := freetype.Pt(paddingLeft, clip.Max.Y-paddingBottom)
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
