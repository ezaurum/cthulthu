package context

import (
	"errors"
	"github.com/ezaurum/cthulthu/errres"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Request struct {
	txRequest
	sessionRequest
	resWriter
	Resource     interface{}
	ResourceName string
	Grant        string
	echoContext  echo.Context
	Notify       func(string, interface{})
}

func New(ctx echo.Context) *Request {
	return &Request{
		echoContext: ctx,
	}
}

func NewWithDB(ctx echo.Context, writer *gorm.DB, reader *gorm.DB) *Request {
	return &Request{
		txRequest: txRequest{
			writeDB: writer,
			readDB:  reader,
		},
		echoContext: ctx,
	}
}

func (r *Request) Context() echo.Context {
	return r.echoContext
}

func (r *Request) HandlerError() error {
	var httpError *errres.HttpError
	if errors.As(r.Error, &httpError) {
		return r.echoContext.JSON(httpError.Code, httpError)
	} else {
		r.echoContext.Error(r.Error)
		return nil
	}
}

func (r *Request) SendResponse() error {
	return r.Complete(r.echoContext)
}

func (r *Request) Request() *http.Request {
	return r.Context().Request()
}
func (r *Request) Response() http.ResponseWriter {
	return r.Context().Response()
}
