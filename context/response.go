package context

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ResponseWriter interface {
	Complete(c echo.Context) error
}

type resWriter struct {
	Result     interface{}
	ResultType string
}

var _ ResponseWriter = &resWriter{}

func (r *resWriter) Complete(c echo.Context) error {
	// 에러 없는 경우
	switch r.ResultType {
	case echo.MIMEApplicationJSON:
		fallthrough
	case echo.MIMEApplicationJSONCharsetUTF8:
		fallthrough
	case "json":
		fallthrough
	default:
		return c.JSON(http.StatusOK, r.Result)
	case "html":
		//todo
		return c.HTML(http.StatusOK, "todo")
	}
}
