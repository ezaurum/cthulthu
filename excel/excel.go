package excel

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tealeg/xlsx"
	"net/http"
	"time"
)

func Serve(c echo.Context, buffer *bytes.Buffer, filename string) {

	response := c.Response()
	request := c.Request()
	response.Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", "attatchment", filename))

	http.ServeContent(response, request, filename, time.Now(), bytes.NewReader(buffer.Bytes()))
}

func NewSheet(file *xlsx.File, name string) (*xlsx.Sheet, error) {
	if sheet, err := file.AddSheet(name); nil != err {
		return nil, err
	} else {
		return sheet, nil
	}
}

func NewSheetWithTitle(file *xlsx.File, name string, titles []string) (*xlsx.Sheet, error) {
	if sheet, err := NewSheet(file, name); nil != err {
		return nil, err
	} else {
		row := sheet.AddRow()
		for _, t := range titles {
			cell := row.AddCell()
			cell.Value = t
		}

		return sheet, nil
	}
}
