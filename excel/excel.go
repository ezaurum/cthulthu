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
	reader := bytes.NewReader(buffer.Bytes())

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", "attatchment", filename))

	http.ServeContent(c.Response(), c.Request(), "attendant-point.xlsx", time.Now(), reader)
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
