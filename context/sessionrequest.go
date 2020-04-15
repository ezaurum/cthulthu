package context

import (
	"fmt"
	"github.com/ezaurum/cthulthu/cookie"
	"github.com/ezaurum/cthulthu/session"
)

type sessionRequest struct {
	Cookie  cookie.Jar
	Session *session.Session
}

func (r *sessionRequest) LoadSession(scn string, maxAge int) {
	if ss, err := session.FromCookie(r.Cookie.Get(scn), maxAge); nil == err {
		r.Session = ss
	} else {
		//todo 세션 생성 에러 처리
		switch err {
		case session.ErrSessionEmpty:
			r.Session = ss
		}
		fmt.Println(err)
	}
}

func (r *sessionRequest) SaveSession(scn string, clientCookieName string, domain string) {
	if httpCookie, clientCookie, err := r.Session.ToCookie(scn, clientCookieName, domain); nil != err {
		//todo 에러 정리
		fmt.Println(err)
	} else {
		r.Cookie.Set(httpCookie)
		r.Cookie.Set(clientCookie)
	}
}
