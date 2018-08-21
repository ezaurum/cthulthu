package barcode

import (
	"testing"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/stretchr/testify/assert"
	"github.com/ezaurum/cthulthu/paint"
)

func TestMakeBarCode(t *testing.T) {
	generate := snowflake.New(0).Generate()
	e, b := MakeMMSBarCodeFile(generate, generate+".jpg", nil)
	assert.True(t, b)
	assert.Nil(t, e)
}

func TestMakeBarCode2(t *testing.T) {
	generate := snowflake.New(0).Generate()
	jpeg := paint.LoadJPEG("dd.jpg")
	e, b := MakeMMSBarCodeFile(generate, generate+".jpg", jpeg)
	assert.True(t, b)
	assert.Nil(t, e)
}
