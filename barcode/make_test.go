package barcode

import (
	"fmt"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/paint"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeBarCode(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	//generate := snowflake.New(0).Generate()
	generate := "180605911050"
	e, b := MakeMMSBarCodeFile(generate, generate+".jpg", nil, nil)
	assert.True(t, b)
	assert.Nil(t, e)
}

func TestMakeBarCode2(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	jpeg := paint.LoadJPEG("wisdom_default.jpg")
	e, b := MakeMMSBarCodeFile(generate, generate+".jpg", jpeg, map[string]interface{}{
		"withCodeString": true,
		"barcodeType":    "128",
	})
	assert.True(t, b)
	assert.Nil(t, e)
}

func TestMakeBarCodeWithString(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	img, err := MakeBarCodeWithString(generate, "128", 100, 10, 32.0, 72.0)

	paint.CreateJPEG(generate+".jpg", img)

	assert.NotNil(t, img)
	assert.Nil(t, err)

}

func TestMakeBarCodeWithString2(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	//generate := snowflake.New(0).Generate()
	generate := "180605911050"
	e, b := MakeMMSBarCodeFile(generate, generate+"TestMakeBarCodeWithString2.jpg", nil, map[string]interface{}{
		"withCodeString": true,
		"barcodeType":    "128",
	})
	assert.True(t, b)
	assert.Nil(t, e)
	MakeMMSBarCodeFile(generate, generate+"TestMakeBarCodeWithString2-25.jpg", nil, map[string]interface{}{
		"withCodeString": true,
		"barcodeType":    "25",
	})
	MakeMMSBarCodeFile(generate, generate+"TestMakeBarCodeWithString2-39.jpg", nil, map[string]interface{}{
		"withCodeString": true,
		"barcodeType":    "39",
	})
	MakeMMSBarCodeFile(generate[0:8], generate+"TestMakeBarCodeWithString2-128.jpg", nil, map[string]interface{}{
		"withCodeString": true,
		"paddingLeft":    200,
		"paddingBottom":  20,
		"fontSize":       48.0,
		"barcodeType":    "128",
	})

}

func TestMakeBarCodeWithString3(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	jpeg := paint.LoadJPEG("dd.jpg")
	e, b := MakeMMSBarCodeFile(generate[0:12], generate+"TestMakeBarCodeWithString3.jpg", jpeg, map[string]interface{}{
		"withCodeString": true,
		"paddingLeft":    140,
		"paddingBottom":  20,
		"fontSize":       48.0,
		"barcodeType":    "128",
	})
	e, b = MakeMMSBarCodeFile(fmt.Sprintf("%X", generate), generate+"TestMakeBarCodeWithString3-2.jpg", jpeg, map[string]interface{}{
		"withCodeString": true,
		"paddingLeft":    200,
		"paddingBottom":  20,
		"fontSize":       48.0,
		"barcodeType":    "128",
	})
	assert.True(t, b)
	assert.Nil(t, e)
}
