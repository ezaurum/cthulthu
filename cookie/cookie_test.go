package cookie

import (
	"github.com/ezaurum/cthulthu/test"
	"github.com/ezaurum/session"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	ct "github.com/ezaurum/cthulthu"
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

//TODO create
//TODO activate
//TODO load auth - 여기서 valid, invalid를 체크해야 한다. has auth info 정도가 좋을까?
//TODO create auth - 실제로 auth 를 생성하는 건 아니다. 그러면 여기서는 create 가 아니라 first? no auth info?
//TODO invalid cookie
//TODO no session 이지 invalid cookie 가 아니라
