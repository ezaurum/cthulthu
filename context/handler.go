package context

import (
	"github.com/ezaurum/cthulthu/cookie"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HandlerFuncResource struct {
	Resource
	Method      string
	Path        string
	HandlerFunc []RequestHandlerFunc
}

type RequestHandlerFunc func(c *Request) error

type ResponseWriter interface {
	Complete(c echo.Context) error
}

func HandlerError(r *Request, c echo.Context) {
	c.JSON(http.StatusInternalServerError, r.Error)
}

// 세션 사용, 트랜잭션 사용
func DefaultHandler(ctx Context, logicArray ...RequestHandlerFunc) func(c echo.Context) error {
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

		// 트랜잭션 완료
		if txErr := r.CompleteTx(); nil != txErr {
			// 트랜잭션이 잘못되면 일단 에러 표시
			r.Error = txErr
			HandlerError(r, c)
		}

		// 결과 전송
		// 에러 있는 경우
		if nil != err {
			HandlerError(r, c)
			return nil
		}

		// 세션 쓰기 - 쿠키로
		clientCookieName := "ss-cck"
		r.SaveSession(scn, clientCookieName, domain)
		// 쿠키 쓰기
		r.Cookie.Write()

		return r.Complete(c)
	}
}

// 세션 사용, 트랜잭션 미 사용
func ReadOnlyHandler(ctx Context, logicArray ...RequestHandlerFunc) func(c echo.Context) error {
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
			HandlerError(r, c)
			return nil
		}

		// 세션 쓰기 - 쿠키로
		clientCookieName := "ss-cck"
		r.SaveSession(scn, clientCookieName, domain)
		// 쿠키 쓰기
		r.Cookie.Write()

		return r.Complete(c)
	}
}

func newRequest(c echo.Context, ctx Context) *Request {
	repo := ctx.Repository()
	writer := repo.Writer()
	reader := repo.Reader()
	// 컨텍스트 생성
	r := NewWithDB(c, writer, reader)
	return r
}
