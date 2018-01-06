package authorizer

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlocked(t *testing.T) {
	r := getDefault(t, TestIDToken{
		Token:        "testToken",
		UserID:       "account",
		UserPassword: "WTF",
	}, TestIdentity{},"test/model.conf", "test/all-access.csv")

	r.GET("/blocked", func(c *gin.Context) {
		assert.Fail(t, "blocked is not runnable")
	})

	requested := false
	r.GET("/ok", func(c *gin.Context) {
		requested = true
	})

	test.GetRequest(r, "/blocked")
	test.GetRequest(r, "/ok")

	assert.True(t, requested, "ok is not requested.")
}

func TestRoleAccess(t *testing.T) {
	token := TestIDToken{
		Token:        "testToken",
		UserID:       "account",
		UserPassword: "WTF",
	}

	identity := TestIdentity{
		IdentityRole:"SuperAdmin",
	}

	r := getDefault(t, token,identity, "test/model.conf", "test/role-access.csv")

	count := 0
	r.GET("/blocked", func(c *gin.Context) {
		count++
	})

	// 강제 로그인
	r.GET("/login", func(c *gin.Context) {
		ac := authenticator.GetAuthenticator(c)
		ac.Authenticate(c, session.GetSession(c), token)
	})

	client := test.HttpClient{}
	client.GetRequest(r, "/blocked")
	assert.Equal(t, 0, count)

	client.GetRequest(r, "/login")
	client.GetRequest(r, "/blocked")
	assert.Equal(t, 1, count)
}

func getDefault(t *testing.T, token authenticator.IDToken,
	identity authenticator.Identity, params ...interface{}) *gin.Engine {
	r := gin.New()

	ac := authenticator.NewMem(0, 100)

	ac.LoadIDToken = func( s string) (authenticator.IDToken, bool) {
		assert.Equal(t, s, token.TokenString())
		return token, true
	}
	ac.LoadIdentity = func( tk authenticator.IDToken) (authenticator.Identity, bool) {
		return identity, true
	}

	assert.NotNil(t, ac.LoadIDToken)
	assert.NotNil(t, ac.LoadIdentity)

	r.Use(ac.Handler())

	authorizer := New(params...)
	auth := AuthorizeMiddleware{
		authorizer: authorizer,
	}

	r.Use(auth.Handler())

	return r
}

type TestIDToken struct {
	id           int64
	UserID       string
	UserPassword string
	isPersisted  bool
	expires      time.Time
	Token        string
	IdentityID int64
}

func (l TestIDToken) TokenString() string {
	return l.Token
}

func (l TestIDToken) IsPersisted() bool {
	return l.isPersisted
}
func (l TestIDToken) IdentityKey() int64 {
	return l.IdentityID
}

func (l TestIDToken) IsExpired() bool {
	return time.Now().After(l.expires)
}

type TestIdentity struct {
	IdentityRole string
}

func (i TestIdentity) Role() string {
	return i.IdentityRole
}

