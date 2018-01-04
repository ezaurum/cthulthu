package authenticator

import (
	ct "github.com/ezaurum/cthulthu"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/session/stores/memstore"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	sessionIDCookieName        = "session-id-CTHULTHU"
	persistedIDTokenCookieName = "persisted-id-token-CTHULTHU"
)

func Init(r *gin.Engine) ct.GinMiddleware {
	ca := Default()
	r.Use(ca.Handler())

	return &ca
}

type IDLoader func(*gin.Context, string) (IDToken, bool)

type cookieAuthenticator struct {
	store                      session.Store
	MaxAge                     int
	sessionIDCookieName        string
	persistedIDTokenCookieName string
	LoadIDToken                IDLoader
}

func (ca cookieAuthenticator) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		session, needSession := ca.findSession(c)

		if needSession {
			session = ca.createSession(c)
		} else if HasIDToken(session) {
			//activate
			ca.activateSession(c, session)
			return
		}

		// remember me
		IDToken, needAuthentication := ca.findIDToken(c, session)
		if needAuthentication {
			ca.Authenticate(c, session, IDToken)
		}

		//activate
		ca.activateSession(c, session)
	}
}

func (ca cookieAuthenticator) createSession(c *gin.Context) session.Session {
	// created
	session := ca.store.GetNew()

	return session
}

func (ca cookieAuthenticator) PersistIDToken(c *gin.Context, session session.Session, idToken IDToken) {
	c.SetCookie(ca.persistedIDTokenCookieName, idToken.TokenString(),
		365*24*60*60*10, "", "", false, true)
}

func (ca cookieAuthenticator) findIDToken(c *gin.Context, session session.Session) (IDToken, bool) {

	//TODO 쿠키가 유효한지 체크를 해 봐야지.
	loginCookie, e := c.Cookie(ca.persistedIDTokenCookieName)
	hasCookie := nil == e

	if hasCookie {

		//IDToken 가져오기
		if nil != ca.LoadIDToken {
			idToken, exist := ca.LoadIDToken(c, loginCookie)

			if exist {
				return idToken, true
			}

		}

		// TODO

		//TODO 가져오기가 없으면 무시?

	}

	// TODO 필요하면 DB에다 회원가입을 시키고?

	return nil, false
}

func (ca cookieAuthenticator) Authenticate(c *gin.Context, session session.Session, IDToken IDToken) {
	SetIDToken(session, IDToken)

	if IDToken.IsPersisted() {
		ca.PersistIDToken(c, session, IDToken)
	}
}

func (ca cookieAuthenticator) activateSession(c *gin.Context, s session.Session) {
	//refresh session expires
	session.SetSession(c, s)
	ca.SetSessionIDCookie(c, s)
}

func (ca cookieAuthenticator) findSession(c *gin.Context) (session.Session, bool) {
	sessionIDCookie, e := c.Cookie(ca.sessionIDCookieName)
	isSessionIdCookieExist := nil == e
	var session session.Session
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
	c.SetCookie(ca.sessionIDCookieName, "", -1, "/", "", false, true)
}

func (ca cookieAuthenticator) SetSessionIDCookie(c *gin.Context, session session.Session) {
	c.SetCookie(ca.sessionIDCookieName, session.ID(),
		ca.MaxAge,
		"/", "",
		false, true)
}

func Default() ct.GinMiddleware {
	return NewMem(0, session.DefaultSessionExpires)
}

func NewMem(node int64, expiresInSeconds int) *cookieAuthenticator {
	duration := time.Duration(expiresInSeconds) * time.Second
	k := snowflake.New(node)
	m := memstore.New(k, duration, duration*2)
	middleware := newMiddleware(m)
	middleware.MaxAge = expiresInSeconds
	return middleware
}

func newMiddleware(store session.Store) *cookieAuthenticator {
	return &cookieAuthenticator{
		store:  store,
		MaxAge: session.DefaultSessionExpires,
		persistedIDTokenCookieName: persistedIDTokenCookieName,
		sessionIDCookieName:        sessionIDCookieName,
	}
}

var _ ct.GinMiddleware = &cookieAuthenticator{}
