package barcode

import (
	"testing"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestMakeBarCode(t *testing.T) {
	e, b := MakeMMSBarCodeFile(snowflake.New(0).Generate(), "test.jpg", nil)
	assert.True(t, b)
	assert.Nil(t, e)
}
