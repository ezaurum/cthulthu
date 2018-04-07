package render

import (
	"net/http"
	"testing"

	"github.com/ezaurum/cthulthu/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	r := getDefault()

	givenUrl := "/"

	testString := "Test"

	r.GET(givenUrl, func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{"TestString": testString})
	})

	doc := test.GetStatusOKDoc(r, givenUrl, t)

	assert.Equal(t, testString, doc.Find("p").First().Text())
	assert.Equal(t, "Dashboard", doc.Find("h1").First().Text())
}

func TestLogin(t *testing.T) {
	r := getDefault()

	givenUrl := "/login"

	r.GET(givenUrl, func(c *gin.Context) {
		c.HTML(http.StatusOK, "login/form", nil)
	})

	w := test.GetStatusOKDoc(r, givenUrl, t)

	assert.Equal(t, 1, w.Find("form").Length())
	assert.Equal(t, "로그인", w.Find("title").First().Text())
}

func TestTemplateDebug(t *testing.T) {
	r := getDefault()
	render := Default()
	r.HTMLRender = render
}

// test utils

func getDefault() *gin.Engine {
	r := gin.New()
	MakeRender(r)
	return r
}
