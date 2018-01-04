package render

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
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

	doc := getStatusOKDoc(r, givenUrl, t)

	assert.Equal(t, testString, doc.Find("p").First().Text())
	assert.Equal(t, "Dashboard", doc.Find("h1").First().Text())
}

func TestLogin(t *testing.T) {
	r:=getDefault()

	givenUrl := "/login"

	r.GET(givenUrl, func(c *gin.Context) {
		c.HTML(http.StatusOK, "login/form",nil)
	})

	w := getStatusOKDoc(r, givenUrl, t)

	assert.Equal(t, 1, w.Find("form").Length())
	assert.Equal(t, "로그인", w.Find("title").First().Text())
}

// test utils

func getStatusOKDoc(r *gin.Engine, givenUrl string, t *testing.T) *goquery.Document {
	w := getRequestResult(r, givenUrl)
	assert.Equal(t, http.StatusOK, w.Code)
	doc := getDocument(w, t)
	return doc
}

func getDefault() *gin.Engine {
	r := gin.New()
	render := Default()
	r.HTMLRender = render
	return r
}

func getDocument(w *httptest.ResponseRecorder, t *testing.T) *goquery.Document {
	doc, error := goquery.NewDocumentFromReader(w.Body)
	if nil != error {
		assert.Fail(t, "failed", error)
	}
	return doc
}

func getRequestResult (r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

