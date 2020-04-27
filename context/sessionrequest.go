package context

import (
	"fmt"
	"github.com/ezaurum/cthulthu/cookie"
	"github.com/ezaurum/cthulthu/session"
	"github.com/labstack/echo/v4"
)

type sessionRequest struct {
	Cookie  cookie.Jar
	Session *session.Session
}

func PopulateSessionFromCookie(c echo.Context, r *Request, ctx Application) error {
	// 쿠키 읽기
	r.Cookie = cookie.New(c.Request(), c.Response())

	scn := ctx.SessionCookieName()
	maxAge := ctx.SessionLifeLength()

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
			fmt.Println("load session error " + err.Error())
			return err
		}
	}
}

func (r *sessionRequest) SaveSession(scn string, clientCookieName string, domain string) error {
	if httpCookie, clientCookie, err := r.Session.ToCookie(scn, clientCookieName, domain); nil != err {
		return err
	} else {
		r.Cookie.Set(httpCookie)
		r.Cookie.Set(clientCookie)
		return nil
	}
	return nil
}

func PopulateSessionFromHeader(r *Request, token string) error {
	if fromToken, err := session.FromToken(token); nil != err {
		return err
	} else {
		r.Session = fromToken
		return nil
	}
}

func WriteSessionToCookie(r *Request, ctx Application) error {
	// 세션 쓰기 - 쿠키 jar에 저장
	if err := r.SaveSession(ctx.SessionCookieName(), ctx.PersistedCookieName(), ctx.Domain()); nil != err {
		return err
	}
	// 쿠키 쓰기
	r.Cookie.Write()
	return nil
}
