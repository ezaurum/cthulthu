package request

import (
	"github.com/ezaurum/cthulthu/context"
	"github.com/ezaurum/cthulthu/cookie"
	"github.com/ezaurum/cthulthu/errres"
	"github.com/ezaurum/cthulthu/session"
	"net/http"
)

func HandleContentTypeSetter(r *context.Request) error {
	r.ResultType = r.Request().Header.Get("Content-Type")
	return nil
}

func HandleTokenSession(c *context.Request) error {
	token := c.Request().Header.Get("CTHULTHU-Token")
	if len(token) < 1 {
		return nil
	}
	if fromToken, err := session.FromToken(token); nil == err {
		c.Session = fromToken
		return nil
	} else {
		return errres.BadReq("invalid token", err, token)
	}
}

func HandleWriteCookieSession(r *context.Request) error {
	// 세션 쓰기 - 쿠키 jar에 저장
	if err := r.SaveSession("scn", "pcn", "localhost"); nil != err {
		return err
	}
	return nil
}

func HandleWriteCookie(r *context.Request) error {
	// 쿠키 쓰기
	r.Cookie.Write()
	return nil
}

func HandleCookieMake(r *context.Request) error {
	// 쿠키 읽기
	r.Cookie = cookie.New(r.Request(), r.Response())
	return nil
}

func HandleCookieSessionPopulate(r *context.Request) error {

	// 필요하면 여기를 변경해서 새로 만든다
	scn := "scn"
	maxAge := 3600

	if ss, err := session.FromCookie(r.Cookie.Get(scn)); nil == err {
		ss.Extends()
		r.Session = ss
		return nil
	} else {
		//세션 생성 에러 처리
		switch err {
		case session.ErrSessionEmpty:
			fallthrough
		case session.ErrSessionExpired:
			sss, _ := session.PopulateAnonymous(maxAge)
			r.Session = &sss
			return nil
		default:
			return errres.WrapWithCode(http.StatusBadRequest, "load session error ", err)
		}
	}
}

// 트랜잭션 완료 핸들러
func HandleCompleteTx(r *context.Request) error {
	if txErr := r.CompleteTx(); nil != txErr {
		return errres.WrapWithCode(http.StatusInternalServerError, "tx complete error", txErr)
	}
	return nil
}
