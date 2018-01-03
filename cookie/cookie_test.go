package cookie

import (
	ct "github.com/ezaurum/cthulthu"
	"github.com/ezaurum/cthulthu/test"
	"github.com/ezaurum/session"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
	"strings"
)

func TestCookie(t *testing.T) {

	r := gin.New()
	r.Use(Default().Handler())

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(ct.DefaultSessionContextKey).(session.Session)
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
	r.Use(n.Handler())

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(ct.DefaultSessionContextKey).(session.Session)
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

func TestPersistToken (t *testing.T) {
	r := gin.New()
	middleware := NewMem(0, 1)
	r.Use(middleware.Handler())

	token := ct.LoginIdentity{
		UserID: "test",
		UserPassword: "test",
		Token : "WTF",
	}

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(ct.DefaultSessionContextKey).(session.Session)
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
func TestPersistedTokenLoad (t *testing.T) {
	r := gin.New()
	middleware := NewMem(0, 1)

	token := ct.LoginIdentity{
		UserID: "test",
		UserPassword: "test",
		Token : "WTF",
	}

	middleware.LoadIDToken = func(context *gin.Context, s string) (ct.IDToken, bool) {
		assert.Equal(t, s, token.TokenString())
		return token, true
	}

	assert.NotNil(t, middleware.LoadIDToken)

	r.Use(middleware.Handler())

	r.GET("/", func(c *gin.Context) {
		s := c.MustGet(ct.DefaultSessionContextKey).(session.Session)
		middleware.PersistIDToken(c, s, token)
		c.String(http.StatusOK, token.TokenString())
	})

	r.GET("/1", func(c *gin.Context) {
		s := c.MustGet(ct.DefaultSessionContextKey).(session.Session)
		t := ct.GetIDToken(s)
		c.String(http.StatusOK, t.TokenString())
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

//TODO create
//TODO activate
//TODO load auth - 여기서 valid, invalid를 체크해야 한다. has auth info 정도가 좋을까?
//TODO create auth - 실제로 auth 를 생성하는 건 아니다. 그러면 여기서는 create 가 아니라 first? no auth info?
//TODO invalid cookie
//TODO no session 이지 invalid cookie 가 아니라
