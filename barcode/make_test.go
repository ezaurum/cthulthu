package barcode

import (
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/paint"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeBarCode(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	e, b := MakeMMSBarCodeFile(generate, generate+".jpg", nil, false)
	assert.True(t, b)
	assert.Nil(t, e)
}

func TestMakeBarCode2(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	jpeg := paint.LoadJPEG("wisdom_default.jpg")
	e, b := MakeMMSBarCodeFile(generate, generate+".jpg", jpeg, false)
	assert.True(t, b)
	assert.Nil(t, e)
}

func TestMakeBarCodeWithString(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	img, err := MakeBarCodeWithString(generate)

	paint.CreateJPEG(generate+".jpg", img)

	assert.NotNil(t, img)
	assert.Nil(t, err)

}

func TestMakeBarCodeWithString2(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()
	e, b := MakeMMSBarCodeFile(generate, generate+"TestMakeBarCodeWithString2.jpg", nil, true)
	assert.True(t, b)
	assert.Nil(t, e)
}

func TestMakeBarCodeWithString3(t *testing.T) {
	//generate := fmt.Sprintf("%X", snowflake.New(0).GenerateInt64())[3:]
	generate := snowflake.New(0).Generate()[0:5]
	jpeg := paint.LoadJPEG("dd.jpg")
	e, b := MakeMMSBarCodeFile(generate, generate+"TestMakeBarCodeWithString3.jpg", jpeg, true)
	assert.True(t, b)
	assert.Nil(t, e)
}
