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
	generate := snowflake.New(0).Generate()
	e, b := MakeMMSBarCodeFile(generate, generate+".jpg", nil, false, "128")
	assert.True(t, b)
	assert.Nil(t, e)
}

func TestMakeBarCode2(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	jpeg := paint.LoadJPEG("wisdom_default.jpg")
	e, b := MakeMMSBarCodeFile(generate, generate+".jpg", jpeg, false, "128")
	assert.True(t, b)
	assert.Nil(t, e)
}

func TestMakeBarCodeWithString(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	img, err := MakeBarCodeWithString(generate, "128")

	paint.CreateJPEG(generate+".jpg", img)

	assert.NotNil(t, img)
	assert.Nil(t, err)

}

func TestMakeBarCodeWithString2(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	e, b := MakeMMSBarCodeFile(generate, generate+"TestMakeBarCodeWithString2.jpg", nil, true, "128")
	assert.True(t, b)
	assert.Nil(t, e)
	 MakeMMSBarCodeFile(generate, generate+"TestMakeBarCodeWithString2-25.jpg", nil, true, "25")
	 MakeMMSBarCodeFile(generate, generate+"TestMakeBarCodeWithString2-39.jpg", nil, true, "39")
	 MakeMMSBarCodeFile(generate[0:8], generate+"TestMakeBarCodeWithString2-128.jpg", nil, true, "128")
}

func TestMakeBarCodeWithString3(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	jpeg := paint.LoadJPEG("dd.jpg")
	e, b := MakeMMSBarCodeFile(generate[0:10], generate+"TestMakeBarCodeWithString3.jpg", jpeg, true, "128")
	e, b = MakeMMSBarCodeFile(fmt.Sprintf("%X", generate), generate+"TestMakeBarCodeWithString3-2.jpg", jpeg, true, "128")
	assert.True(t, b)
	assert.Nil(t, e)
}
