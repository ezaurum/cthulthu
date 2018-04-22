package authenticator

import (
	"fmt"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCookie(t *testing.T) {

	r := gin.New()
	authenticator := getDefault()
	r.Use(authenticator.Handler())

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(session.DefaultSessionContextKey).(session.Session)
		c.String(http.StatusOK, s.ID())
	})

	client := test.HttpClient{}

	// first request
	w := client.GetRequest(r, "/")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), test.GetFirstCookieValue(client.Cookies))

	// second request with cookie
	w0 := client.GetRequest(r, "/")

	assert.Equal(t, http.StatusOK, w0.Code)
	assert.Equal(t, w0.Body.String(), w.Body.String())
}
func TestExpire(t *testing.T) {

	r := gin.New()
	n := NewMem(0, 1)
	setNilFunctions(n)
	r.Use(n.Handler())

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(session.DefaultSessionContextKey).(session.Session)
		c.String(http.StatusOK, s.ID())
	})

	client := test.HttpClient{}

	// first request
	w := client.GetRequest(r, "/")

	time.Sleep(time.Second * 1)

	// second request with cookie
	w0 := client.GetRequest(r, "/")

	// should not be equal
	assert.NotEqual(t, w0.Body.String(), w.Body.String())
}

func TestPersistToken(t *testing.T) {
	r := gin.New()
	middleware := NewMem(0, 1)
	setNilFunctions(middleware)
	r.Use(middleware.Handler())

	token := FormIDToken{
		UserID:       "test",
		UserPassword: "test",
		Token:        "WTF",
	}

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(session.DefaultSessionContextKey).(session.Session)
		middleware.PersistIDToken(c, s, token)
		c.String(http.StatusOK, token.TokenString())
	})

	client := test.HttpClient{}

	// first request
	client.GetRequest(r, "/")

	assert.True(t, len(client.Cookies) > 1)

	e := false
	for _, c := range client.Cookies {
		e = strings.Contains(c, middleware.persistedIDTokenCookieName)
		if e {
			break
		}
	}

	assert.True(t, e)
}

func TestPersistedTokenLoad(t *testing.T) {
	r := gin.New()

	token := FormIDToken{
		UserID:       "test",
		UserPassword: "test",
		Token:        "WTF",
	}

	identity := TestIdentity{
		IdentityRole: "Role",
	}

	middleware := getTestAuthenticator(t, token, identity, 1)
	middleware.SetActions(func(s string) (IDToken, bool) {
		return token, true
	}, func(token IDToken) (Identity, bool) {
		return identity, true
	}, func(token IDToken) IDToken {
		return nil
	})

	r.Use(middleware.Handler())

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(session.DefaultSessionContextKey).(session.Session)
		middleware.PersistIDToken(c, s, token)
		c.String(http.StatusOK, token.TokenString())
	})

	r.GET("/1", func(c *gin.Context) {
		s := session.GetSession(c)
		tk, _ := s.Get(IDTokenSessionKey)

		i, _ := s.Get(IdentitySessionKey)
		c.String(http.StatusOK, tk.(IDToken).TokenString())
		assert.Equal(t, identity.Role(), i.(Identity).Role())
	})

	client := test.HttpClient{}

	// first request
	w0 := client.GetRequest(r, "/")
	w0Cookie := test.GetFirstCookieValue(client.Cookies)

	// expired
	time.Sleep(time.Second)

	w1 := client.GetRequest(r, "/1")
	w1Cookie := test.GetFirstCookieValue(client.Cookies)

	assert.Equal(t, w0.Body.String(), w1.Body.String())

	assert.NotEqual(t, w0Cookie, w1Cookie)
}
func TestIdentityRole(t *testing.T) {
	r := gin.New()

	token := FormIDToken{
		UserID:       "test",
		UserPassword: "test",
		Token:        "WTF",
	}

	identity := TestIdentity{
		IdentityRole: "Role",
	}

	middleware := getTestAuthenticator(t, token, identity, 10)
	setNilFunctions(middleware)
	middleware.LoadIdentity = func(tkn IDToken) (Identity, bool) {
		if tkn.TokenString() == token.TokenString() {
			return identity, true
		}
		return nil, false
	}

	r.Use(middleware.Handler())

	r.GET("/", func(c *gin.Context) {
		ac := GetAuthenticator(c)
		ac.Authenticate(c, session.GetSession(c), token)

		fmt.Println(session.GetSession(c))
	})

	r.GET("/1", func(c *gin.Context) {
		fmt.Println(session.GetSession(c))
		id := GetIdentity(session.GetSession(c))
		assert.Equal(t, identity.Role(), id.Role())
	})

	client := test.HttpClient{}

	// first request
	client.GetRequest(r, "/")
	client.GetRequest(r, "/1")
}

func getTestAuthenticator(t *testing.T, token IDToken, identity Identity, expiresInSecond int) *cookieAuthenticator {

	middleware := NewMem(0, expiresInSecond)
	setNilFunctions(middleware)

	middleware.LoadIDToken = func(s string) (IDToken, bool) {
		assert.Equal(t, s, token.TokenString())
		return token, true
	}

	middleware.LoadIdentity = func(token IDToken) (Identity, bool) {
		assert.NotNil(t, token)
		assert.NotNil(t, identity)
		return identity, true
	}

	assert.NotNil(t, middleware.LoadIDToken)
	return middleware
}

type FormIDToken struct {
	id           int64
	UserID       string
	UserPassword string
	isPersisted  bool
	expires      time.Time
	Token        string
	IdentityID   int64
}

func (l FormIDToken) IdentityKey() int64 {
	return l.IdentityID
}

func (l FormIDToken) TokenString() string {
	return l.Token
}

func (l FormIDToken) IsPersisted() bool {
	return l.isPersisted
}

func (l FormIDToken) IsExpired() bool {
	return time.Now().After(l.expires)
}

type TestIdentity struct {
	IdentityRole string
}

func (i TestIdentity) Role() string {
	return i.IdentityRole
}

func getDefault() *cookieAuthenticator {
	authenticator := Default()
	setNilFunctions(authenticator)
	return authenticator
}
func setNilFunctions(authenticator Authenticator) {
	authenticator.SetActions(func(s string) (IDToken, bool) {
		return nil, false
	}, func(token IDToken) (Identity, bool) {
		return nil, false
	}, func(token IDToken) IDToken {
		return nil
	})
}

//TODO create
//TODO activate
//TODO load auth - 여기서 valid, invalid를 체크해야 한다. has auth info 정도가 좋을까?
//TODO create auth - 실제로 auth 를 생성하는 건 아니다. 그러면 여기서는 create 가 아니라 first? no auth info?
//TODO invalid cookie
//TODO no session 이지 invalid cookie 가 아니라
