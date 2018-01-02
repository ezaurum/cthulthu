package test

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"net/url"
	"io"
	"net/http"
	"bytes"
	"encoding/json"
)

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

	req, _ := http.NewRequest("POST", url,b)
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

