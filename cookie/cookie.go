package cookie

import (
	ct "github.com/ezaurum/cthulthu"
	ezs "github.com/ezaurum/session"
	"github.com/ezaurum/session/generators/snowflake"
	"github.com/ezaurum/session/stores/memstore"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	sessionIDCookieName    = "session-id-CTHULTHU"
	persistedIDTokenCookieName = "persisted-id-token-CTHULTHU"
)

type IDLoader func(*gin.Context, string) (ct.IDToken, bool)

type cookieAuthenticator struct {
	store  ezs.Store
	MaxAge int
	sessionIDCookieName string
	persistedIDTokenCookieName string
	LoadIDToken IDLoader
}

func (ca cookieAuthenticator) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		session, needSession := ca.findSession(c)

		if needSession {
			session = ca.createSession(c)
		} else if ct.HasIDToken(session) {
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

func (ca cookieAuthenticator) createSession(c *gin.Context) ezs.Session {
	// created
	session := ca.store.GetNew()

	return session
}

func (ca cookieAuthenticator) PersistIDToken(c *gin.Context, session ezs.Session, idToken ct.IDToken) {
	c.SetCookie(ca.persistedIDTokenCookieName, idToken.TokenString(),
		365*24*60*60*10, "","", false, true)
}

func (ca cookieAuthenticator) findIDToken(c *gin.Context, session ezs.Session) (ct.IDToken, bool) {

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

		//TODO 가져오기가 없으면 무시?

	}

	// TODO 필요하면 DB에다 회원가입을 시키고?

	return nil, false
}

func (ca cookieAuthenticator) Authenticate(c *gin.Context, session ezs.Session, IDToken ct.IDToken) {
	ct.SetIDToken(session, IDToken)

	if IDToken.IsPersisted() {
		ca.PersistIDToken(c, session, IDToken)
	}
}

func (ca cookieAuthenticator) activateSession(c *gin.Context, session ezs.Session) {
	//refresh session expires
	ct.SetSession(c, session)
	ca.SetSessionIDCookie(c, session)
}

func (ca cookieAuthenticator) findSession(c *gin.Context) (ezs.Session, bool) {
	sessionIDCookie, e := c.Cookie(ca.sessionIDCookieName)
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
	c.SetCookie(ca.sessionIDCookieName, "", -1, "/", "", false, true)
}

func (ca cookieAuthenticator) SetSessionIDCookie(c *gin.Context, session ezs.Session) {
	c.SetCookie(ca.sessionIDCookieName, session.ID(),
		ca.MaxAge,
		"/", "",
		false, true)
}

func Default() ct.GinMiddleware {
	return NewMem(0, ct.DefaultSessionExpires)
}

func NewMem(node int64, expiresInSeconds int) *cookieAuthenticator {
	duration := time.Duration(expiresInSeconds) * time.Second
	k := snowflake.New(node)
	m := memstore.New(k, duration, duration*2)
	middleware := newMiddleware(m)
	middleware.MaxAge = expiresInSeconds
	return middleware
}

func newMiddleware(store ezs.Store) *cookieAuthenticator {
	return &cookieAuthenticator{
		store:  store,
		MaxAge: ct.DefaultSessionExpires,
		persistedIDTokenCookieName:persistedIDTokenCookieName,
		sessionIDCookieName:sessionIDCookieName,
	}
}

var _ ct.GinMiddleware = &cookieAuthenticator{}
