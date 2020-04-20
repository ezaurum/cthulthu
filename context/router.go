package context

import (
	"net/http"
)

type Router interface {
	Handlers() []HandlerFuncResource
	GET(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	POST(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	PUT(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	DELETE(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	PATCH(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	TRACE(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	OPTION(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	HEAD(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	CONNECT(path string, handlerFunc RequestHandlerFunc, options ...interface{})
}

type router struct {
	handlers []HandlerFuncResource
}

var _ Router = &router{}

func (a *router) appendAndAssign(path string, handlerFunc RequestHandlerFunc, method string) {
	a.handlers = append(a.handlers, HandlerFuncResource{
		Resource: Resource{
			Name:         "",
			ResourceType: HandlerFuncResourceType,
		},
		Method:      method,
		Path:        path,
		HandlerFunc: handlerFunc,
	})
}

func (a *router) PUT(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodPut)
}

func (a *router) DELETE(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodDelete)
}

func (a *router) PATCH(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodPatch)
}

func (a *router) TRACE(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodTrace)
}

func (a *router) OPTION(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodOptions)
}

func (a *router) HEAD(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodHead)
}

func (a *router) CONNECT(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodConnect)
}

func (a *router) GET(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodGet)
}
func (a *router) POST(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.appendAndAssign(path, handlerFunc, http.MethodPost)
}

func (a *router) Handlers() []HandlerFuncResource {
	return a.handlers
}
