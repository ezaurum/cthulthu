package context

import (
	"github.com/ezaurum/cthulthu/cookie"
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

		// 쿠키 읽기
		r.Cookie = cookie.New(c.Request(), c.Response())

		scn := ctx.SessionCookieName()
		domain := ctx.Domain()
		maxAge := ctx.SessionLifeLength()
		r.LoadSession(scn, maxAge)

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

		// 세션 쓰기 - 쿠키로
		r.SaveSession(scn, ctx.PersistedCookieName() , domain)
		// 쿠키 쓰기
		r.Cookie.Write()

		return r.Response()
	}
}

// 세션 사용, 트랜잭션 미 사용
func ReadOnlyHandler(ctx Application, logicArray ...RequestHandlerFunc) func(c echo.Context) error {
	return func(c echo.Context) error {
		r := newRequest(c, ctx)
		r.ResultType = c.Request().Header.Get("Content-Type")

		// 쿠키 읽기
		r.Cookie = cookie.New(c.Request(), c.Response())

		scn := "session-cookie-name"
		domain := "localhost"
		// 세션 읽기
		//todo
		maxAge := 3600
		r.LoadSession(scn, maxAge)

		//todo 권한 처리
		err := r.RunAll(logicArray)

		// 결과 전송
		// 에러 있는 경우
		if nil != err {
			return r.HandlerError()
		}

		// 세션 쓰기 - 쿠키로
		clientCookieName := "ss-cck"
		r.SaveSession(scn, clientCookieName, domain)
		// 쿠키 쓰기
		r.Cookie.Write()

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
