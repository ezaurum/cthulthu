package authenticator

import (
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/session/memstore"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

const (
	SessionIDCookieName        = "session-id-CTHULTHU"
	PersistedIDTokenCookieName = "persisted-id-token-CTHULTHU"
	DefaultSessionContextKey   = "session context key tekeli-li tekeli-li"
	DefaultSessionExpires      = 60 * 15
)

type CookieAuthenticator struct {
	store                      session.Store
	MaxAge                     int
	sessionIDCookieName        string
	persistedIDTokenCookieName string
	LoadIDToken                IDTokenLoader
	LoadIdentity               IDLoader
	PersistToken               TokenSaver
}

func (ca *CookieAuthenticator) SetActions(loadIDToken IDTokenLoader,
	loadIdentity IDLoader, tokenSaver TokenSaver) {
	ca.LoadIdentity = loadIdentity
	ca.LoadIDToken = loadIDToken
	ca.PersistToken = tokenSaver
}

func (ca *CookieAuthenticator) Handler() gin.HandlerFunc {

	if nil == ca.LoadIDToken {
		panic("LoadIdToken is nil")
	}

	if nil == ca.LoadIdentity {
		panic("LoadIdentity is nil")
	}

	if nil == ca.PersistToken {
		panic("PersistToken is nil")
	}

	return func(c *gin.Context) {

		SetAuthenticator(c, ca)

		session, needSession := ca.findSession(c)

		if needSession {
			session = ca.createSession(c)
		} else if IsAuthenticated(session) {
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

func (ca CookieAuthenticator) createSession(c *gin.Context) session.Session {
	// created
	session := ca.store.GetNew()

	return session
}

func (ca CookieAuthenticator) PersistIDToken(c *gin.Context, session session.Session, idToken IDToken) {
	_ = ca.PersistToken(idToken)
	c.SetCookie(ca.persistedIDTokenCookieName, idToken.TokenString(),
		//tODO max age from token
		365*24*60*60*99, "", "", false, true)
}

func (ca CookieAuthenticator) findIDToken(c *gin.Context, session session.Session) (IDToken, bool) {

	loginCookie, e := c.Cookie(ca.persistedIDTokenCookieName)
	hasCookie := nil == e

	if !hasCookie {
		return nil, false
	}

	idToken, exist := ca.LoadIDToken(loginCookie)
	if exist {
		return idToken, true
	}

	return nil, false
}

func (ca CookieAuthenticator) Authenticate(c *gin.Context, session session.Session, idToken IDToken) {

	identity, b := ca.LoadIdentity(idToken)

	if !b {
		//TODO panic("Not exist identity")
		log.Println("Not exist identity")
		//일단 제거
		ca.ClearSessionIDCookie(c)
		return
	}

	SetIdentity(session, identity)
	SetIDToken(session, idToken)
	session.Save()

}

func (ca *CookieAuthenticator) activateSession(c *gin.Context, s session.Session) {

	//refresh session expires
	SetSession(c, s)
	ca.SetSessionIDCookie(c, s)
}

func (ca CookieAuthenticator) findSession(c *gin.Context) (session.Session, bool) {
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

func (ca CookieAuthenticator) ClearSessionIDCookie(c *gin.Context) {
	c.SetCookie(ca.sessionIDCookieName, "", -1, "/", "", false, true)
}

func (ca CookieAuthenticator) SetSessionIDCookie(c *gin.Context, session session.Session) {
	c.SetCookie(ca.sessionIDCookieName, session.ID(),
		ca.MaxAge,
		"/", "",
		false, true)
}

func Default() *CookieAuthenticator {
	return NewMem(0, DefaultSessionExpires)
}

func NewMem(node int64, expiresInSeconds int) *CookieAuthenticator {
	duration := time.Duration(expiresInSeconds) * time.Second
	k := snowflake.New(node)
	m := memstore.New(k, duration, duration*2)
	middleware := newMiddleware(m)
	middleware.MaxAge = expiresInSeconds
	return middleware
}

func newMiddleware(store session.Store) *CookieAuthenticator {
	return &CookieAuthenticator{
		store:  store,
		MaxAge: DefaultSessionExpires,
		persistedIDTokenCookieName: PersistedIDTokenCookieName,
		sessionIDCookieName:        SessionIDCookieName,
	}
}

var _ Authenticator = &CookieAuthenticator{}

func GetSession(c *gin.Context) session.Session {
	return c.MustGet(DefaultSessionContextKey).(session.Session)
}

func SetSession(c *gin.Context, s session.Session) {
	c.Set(DefaultSessionContextKey, s)
}
