package context

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type Request struct {
	txRequest
	sessionRequest
	resWriter
	Resource     interface{}
	ResourceName string
	Error        error
	Grant        string
	echoContext  echo.Context
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

func (r *Request) HandlerError(c echo.Context) {
	c.Error(r.Error)
}
