package request_test

import (
	"encoding/json"
	"errors"
	"github.com/ezaurum/cthulthu/context"
	"github.com/ezaurum/cthulthu/context/request"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/errres"
	"github.com/ezaurum/cthulthu/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultHandler(t *testing.T) {

	// 테스트 서버
	e := echo.New()

	// 테스트 DB
	repo, closeFunc := test.NewRepo(t, "handler")
	ctx := context.App()
	ctx.SetRepository(repo)
	defer closeFunc()

	// 테스트 핸들러
	handler := request.DefaultHandler(ctx, func(c *context.Request) error {
		c.ResultType = echo.MIMEApplicationJSONCharsetUTF8
		c.Result = "OK"
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	err := handler(c)
	if assert.NoError(t, err) {
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get("Content-Type"))
		assert.Equal(t, "\"OK\"\n", rec.Body.String())
	}
}

func TestTxHandler(t *testing.T) {

	// 테스트 서버
	e := echo.New()

	// 테스트 db
	repo, closeFunc := test.NewRepo(t, "handler")
	ctx := context.App()
	ctx.SetRepository(repo)
	defer closeFunc()

	type reqJSON struct {
		database.Model
		Action string `json:"action"`
	}

	// 테이블 생성
	repo.Writer().AutoMigrate(&reqJSON{})

	// 테스트 핸들러
	handler := request.DefaultHandler(ctx, func(c *context.Request) error {
		var reqJ reqJSON
		if err := c.Bind(&reqJ); nil != err {
			return errres.BadRequest("bind error", err)
		}
		c.Tx().Create(&reqJSON{
			Action: "test0",
		})
		c.Tx().Create(&reqJSON{
			Action: "test1",
		})
		c.Tx().Create(&reqJSON{
			Action: "test2",
		})
		c.ResultType = echo.MIMEApplicationJSONCharsetUTF8
		var reqJJ []reqJSON
		c.Tx().Find(&reqJJ)
		c.Result = reqJJ
		if reqJ.Action == "fail" {
			err := errors.New("failed")
			c.SetTxError(err)
			return err
		}

		return nil
	})

	// 성공
	req0 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"action":"success"}`))
	req0.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec0 := httptest.NewRecorder()

	c0 := e.NewContext(req0, rec0)
	err := handler(c0)
	if assert.NoError(t, err) {
		var result []reqJSON
		er := json.Unmarshal(rec0.Body.Bytes(), &result)
		assert.NoError(t, er)
		assert.Equal(t, 3, len(result))
	}

	// 실패
	req1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"action":"fail"}`))
	req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec1 := httptest.NewRecorder()

	c1 := e.NewContext(req1, rec1)
	err1 := handler(c1)
	assert.NotEqual(t, http.StatusOK, rec1.Code)
	assert.NoError(t, err1)

	// 성공
	req2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"action":"success"}`))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec2 := httptest.NewRecorder()

	c2 := e.NewContext(req2, rec2)
	err2 := handler(c2)
	if assert.NoError(t, err2) {
		var result []reqJSON
		er := json.Unmarshal(rec2.Body.Bytes(), &result)
		assert.NoError(t, er)
		assert.Equal(t, 6, len(result))
	}
}
