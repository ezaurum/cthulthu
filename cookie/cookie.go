package cookie

import (
	ct "github.com/ezaurum/cthulthu"
	ezs "github.com/ezaurum/session"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/ezaurum/session/generators/snowflake"
	"github.com/ezaurum/session/stores/memstore"
)
const (
	sessionIDCookieName    = "ca-default-name"
	persistLoginCookieName = "ca-default-remember-me"
)

type cookieAuthenticator struct {
	store ezs.Store
	MaxAge int
}

func (ca cookieAuthenticator) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		session, needSession := ca.findSession(c)

		if needSession {
			session = ca.createSession(c)
		} else if ct.IsAuthenticated(session) {
			//activate
			ca.activateSession(c, session)
			return
		}

		// remember me
		identifier, needAuthentication := ca.findIdentifier(c, session)
		if needAuthentication {
			ca.authenticate(c, session, identifier)
		}

		//activate
		ca.activateSession(c, session)
	}
}

func (ca cookieAuthenticator) createSession(c *gin.Context) ezs.Session {
	// set cookie
	session := ca.store.GetNew()
	// created

	return session
}

func (ca cookieAuthenticator) findIdentifier(c *gin.Context, session ezs.Session) (ct.Identifier, bool) {

	loginCookie, e := c.Cookie(persistLoginCookieName)
	hasCookie := nil == e

	if hasCookie {

		//TODO 쿠키가 유효한지 체크를 해 봐야지.
		if loginCookie == "WTF" {
			return nil, false
		}

		//TODO db identifier 에서 가져오기

		//TODO valid 한 identifier 생성?
		var identifier ct.Identifier
		//TODO 쿠키에다 값 리프레시?
		return identifier, true
	}

	// TODO 필요하면 DB에다 회원가입을 시키고?

	return nil, false
}

func (ca cookieAuthenticator) authenticate(c *gin.Context, session ezs.Session, identifier ct.Identifier) {
	ct.SetIdentifier(session, identifier)
}

func (ca cookieAuthenticator) activateSession(c *gin.Context, session ezs.Session) {
	//TODO refresh session expires
	c.Set(ct.DefaultSessionContextKey, session)
	ca.SetSessionIDCookie(c, session)
}

func (ca cookieAuthenticator) findSession(c *gin.Context) (ezs.Session, bool) {
	sessionIDCookie, e := c.Cookie(sessionIDCookieName)
	isSessionIdCookieExist := nil == e
	var session ezs.Session
	var noSession bool

	if isSessionIdCookieExist {
		//TODO secure
		s, n := ca.store.Get(sessionIDCookie)
		session = s
		noSession = !n

		if noSession {
			// 세션 유효하지 않은 경우, 만료되었거나, 값 조작이거나
			// 해당 쿠키 삭제
			ca.ClearSessionIDCookie(c)
		}
	}
	needSession := !isSessionIdCookieExist || noSession
	return session, needSession
}

func (ca cookieAuthenticator) ClearSessionIDCookie(c *gin.Context) {
	c.SetCookie(sessionIDCookieName, "", -1, "/", "", false, true)
}

func (ca cookieAuthenticator) SetSessionIDCookie(c *gin.Context, session ezs.Session) {
	c.SetCookie(sessionIDCookieName, session.ID(),
		ca.MaxAge,
		"/", "",
		false, true)
}

func Default() ct.GinMiddleware {
	return NewMem(0, ct.DefaultSessionExpires)
}

func NewMem(node int64, expiresInSeconds int) ct.GinMiddleware {
	duration := time.Duration(expiresInSeconds) * time.Second
	k := snowflake.New(node)
	m := memstore.New(k, duration, duration*2)
	middleware := newMiddleware(m).(*cookieAuthenticator)
	middleware.MaxAge = expiresInSeconds
	return middleware
}

func newMiddleware(store ezs.Store) ct.GinMiddleware {
	return &cookieAuthenticator{
		store: store,
		MaxAge: ct.DefaultSessionExpires,
	}
}

var _ ct.GinMiddleware = &cookieAuthenticator{}
