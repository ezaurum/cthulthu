package request

import (
	"github.com/ezaurum/cthulthu/context"
	"github.com/labstack/echo/v4"
)

// 세션 사용, 트랜잭션 사용
func DefaultHandler(ctx context.Application, logicArray ...context.RequestHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := NewRequest(c, ctx)

		if err := r.RunAll(DefaultOpenHandlers); nil != err {
			return r.HandlerError()
		}

		err := r.RunAll(logicArray)

		if err := r.RunAll(DefaultCloseHandlers); nil != err {
			return r.HandlerError()
		}

		// 결과 전송
		// 에러 있는 경우
		if nil != err {
			return r.HandlerError()
		}

		return r.SendResponse()
	}
}

func NewRequest(c echo.Context, ctx context.Application) *context.Request {
	repo := ctx.Repository()
	writer := repo.Writer()
	reader := repo.Reader()
	// 컨텍스트 생성
	r := context.NewWithDB(c, writer, reader)
	r.Notify = ctx.Notify
	return r
}

var DefaultOpenHandlers = []context.RequestHandlerFunc{
	HandleTokenSession,
	HandleCookieMake,
	HandleCookieSessionPopulate,
	HandleContentTypeSetter,
}

var DefaultCloseHandlers = []context.RequestHandlerFunc{
	HandleCompleteTx,
	HandleWriteCookieSession,
	HandleWriteCookie,
}
