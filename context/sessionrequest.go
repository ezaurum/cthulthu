package context

import (
	cookie "github.com/ezaurum/cookie-jar"
	"github.com/ezaurum/cthulthu/session"
	"github.com/labstack/echo/v4"
)

type sessionRequest struct {
	Cookie  cookie.Jar
	Session *session.Session
}

type SessionCookiePopulatorFunc func(c echo.Context, r *Request, ctx Application) error
type SessionCookieWriterFunc func(r *Request, ctx Application) error

func (r *sessionRequest) SaveSession(scn string, clientCookieName string, domain string) error {
	if httpCookie, clientCookie, err := r.Session.ToCookie(scn, clientCookieName, domain); nil != err {
		return err
	} else {
		r.Cookie.Set(httpCookie)
		r.Cookie.Set(clientCookie)
		return nil
	}
}
