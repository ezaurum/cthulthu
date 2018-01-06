package test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func GetFirstCookieValue(cookie []string) string {
	return strings.Split(strings.Split(cookie[0], ";")[0], "=")[1]
}

func PostRequestWithCookie(r *gin.Engine, url string, body io.Reader, cookie []string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Cookie", strings.Join(cookie, ";"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func GetStatusOKDoc(r *gin.Engine, givenUrl string, t *testing.T) *goquery.Document {
	w := GetRequest(r, givenUrl)
	assert.Equal(t, http.StatusOK, w.Code)
	doc := GetDocument(w, t)
	return doc
}

func GetRequest(r *gin.Engine, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func GetRequestWithCookie(r *gin.Engine, url string, cookie []string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", strings.Join(cookie, ";"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func GetCookie(w *httptest.ResponseRecorder) []string {
	cookie := w.HeaderMap["Set-Cookie"]
	return cookie
}

func GetDocument(w *httptest.ResponseRecorder, t *testing.T) *goquery.Document {
	doc, error := goquery.NewDocumentFromReader(w.Body)
	if nil != error {
		assert.Fail(t, "failed", error)
	}
	return doc
}

func IsRedirect(recorder *httptest.ResponseRecorder) bool {
	location := recorder.HeaderMap["Location"]
	return len(location) > 0
}

func GetRedirect(recorder *httptest.ResponseRecorder) string {
	return recorder.HeaderMap["Location"][0]
}
