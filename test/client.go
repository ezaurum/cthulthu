package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// gin 용 테스트 클라이언트. 쿠키 저장 때문에 따로 만들었다

type HttpClient struct {
	Cookies []string
}

func (c *HttpClient) GetRequest(r *gin.Engine, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	if len(c.Cookies) > 0 {
		req.Header.Set("Cookie", strings.Join(c.Cookies, ";"))
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	cookie := GetCookie(w)
	if len(cookie) > 0 {
		c.Cookies = cookie
	}
	return w
}

func (c *HttpClient) PostFormRequest(r *gin.Engine, url string, values url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(c.Cookies) > 0 {
		req.Header.Set("Cookie", strings.Join(c.Cookies, ";"))
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	cookie := GetCookie(w)
	if len(cookie) > 0 {
		c.Cookies = cookie
	}
	return w
}

func (c *HttpClient) PostJsonRequest(r *gin.Engine, url string, value interface{}) *httptest.ResponseRecorder {

	b := new(bytes.Buffer)
	e := json.NewEncoder(b).Encode(value)
	if nil != e {
		panic(e)
	}

	req, _ := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	if len(c.Cookies) > 0 {
		req.Header.Set("Cookie", strings.Join(c.Cookies, ";"))
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	cookie := GetCookie(w)
	if len(cookie) > 0 {
		c.Cookies = cookie
	}
	return w
}

func (c *HttpClient) PostRequest(r *gin.Engine, url string, body io.Reader) *httptest.ResponseRecorder {

	req, _ := http.NewRequest("POST", url, body)
	if len(c.Cookies) > 0 {
		req.Header.Set("Cookie", strings.Join(c.Cookies, ";"))
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	cookie := GetCookie(w)
	if len(cookie) > 0 {
		c.Cookies = cookie
	}
	return w
}
