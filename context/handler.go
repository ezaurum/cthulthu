package context

import (
	"github.com/ezaurum/cthulthu/errres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type HandlerFuncResource struct {
	Resource
	Method      string
	Path        string
	HandlerFunc []RequestHandlerFunc
}

type RequestHandlerFunc func(c *Request) error

// 세션 사용, 트랜잭션 사용
func DefaultHandler(ctx Application, logicArray ...RequestHandlerFunc) func(c echo.Context) error {
	return func(c echo.Context) error {
		r := newRequest(c, ctx)
		r.ResultType = c.Request().Header.Get("Content-Type")

		token := c.Request().Header.Get(ctx.TokenName())
		if len(token) > 0 {
			PopulateSessionFromHeader(r, token)
		} else {
			if err2 := PopulateSessionFromCookie(c, r, ctx); nil != err2 {
				r.Error = err2
				return r.HandlerError()
			}
		}

		err := r.RunAll(logicArray)

		// 트랜잭션 완료
		if txErr := r.CompleteTx(); nil != txErr {
			log.Errorf("transaction error %v", txErr)
			if nil == r.Error {
				r.Error = errres.Wrap("transaction error", txErr)
				return r.HandlerError()
			}
		}

		// 결과 전송
		// 에러 있는 경우
		if nil != err {
			return r.HandlerError()
		}

		if len(token) < 1 {
			if err2 := WriteSessionToCookie(r, ctx); nil != err2 {
				r.Error = err
				return r.HandlerError()
			}
		}

		return r.Response()
	}
}


func newRequest(c echo.Context, ctx Application) *Request {
	repo := ctx.Repository()
	writer := repo.Writer()
	reader := repo.Reader()
	// 컨텍스트 생성
	r := NewWithDB(c, writer, reader)
	r.Notify = ctx.Notify
	return r
}
