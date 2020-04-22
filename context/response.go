package context

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type resWriter struct {
	StatusCode int
	Error      error
	Result     interface{}
	ResultType string
}

var _ ResponseWriter = &resWriter{}

func (r *resWriter) Complete(c echo.Context) error {
	if r.StatusCode != 0 {
		r.StatusCode = http.StatusOK
	}
	// 에러 없는 경우
	switch r.ResultType {
	case echo.MIMEApplicationJSON:
		fallthrough
	case echo.MIMEApplicationJSONCharsetUTF8:
		fallthrough
	case "json":
		fallthrough
	default:
		return c.JSON(r.StatusCode, r.Result)
		//todo 엑셀
		//todo 이미지
	case "html":
		//todo
		return c.HTML(r.StatusCode, "todo")
	}
}
