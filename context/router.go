package context

import (
	"github.com/labstack/echo/v4"
)

type Router interface {
	HandlerSetter
	RouteGroup
	Assign(e *echo.Echo, ctx Context, handlers ...RequestHandlerFunc) error
}

type HandlerSetter interface {
	Handlers() []HandlerFuncResource
	GET(path string, handlerFunc ...RequestHandlerFunc)
	POST(path string, handlerFunc ...RequestHandlerFunc)
	PUT(path string, handlerFunc ...RequestHandlerFunc)
	DELETE(path string, handlerFunc ...RequestHandlerFunc)
	PATCH(path string, handlerFunc ...RequestHandlerFunc)
	TRACE(path string, handlerFunc ...RequestHandlerFunc)
	OPTION(path string, handlerFunc ...RequestHandlerFunc)
	HEAD(path string, handlerFunc ...RequestHandlerFunc)
	CONNECT(path string, handlerFunc ...RequestHandlerFunc)
	AddHandler(path string, method string, handlerFunc ...RequestHandlerFunc)
}

type router struct {
	handlers      []HandlerFuncResource
	parent        Router
	basePath      string
	groupHandlers []RequestHandlerFunc
	children      []Router
}

var _ Router = &router{}
