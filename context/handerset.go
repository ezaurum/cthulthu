package context

import (
	"net/http"
)

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

func (a *router) AddHandler(pathString string, method string, handlerFunc ...RequestHandlerFunc) {
	a.handlers = append(a.handlers, HandlerFuncResource{
		Resource: Resource{
			Name:         "",
			ResourceType: HandlerFuncResourceType,
		},
		Method:      method,
		Path:        pathString,
		HandlerFunc: handlerFunc,
	})
}

func (a *router) PUT(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodPut, handlerFunc...)
}

func (a *router) DELETE(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodDelete, handlerFunc...)
}

func (a *router) PATCH(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodPatch, handlerFunc...)
}

func (a *router) TRACE(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodTrace, handlerFunc...)
}

func (a *router) OPTION(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodOptions, handlerFunc...)
}

func (a *router) HEAD(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodHead, handlerFunc...)
}

func (a *router) CONNECT(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodConnect, handlerFunc...)
}

func (a *router) GET(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodGet, handlerFunc...)
}
func (a *router) POST(path string, handlerFunc ...RequestHandlerFunc) {
	a.AddHandler(path, http.MethodPost, handlerFunc...)
}

func (a *router) Handlers() []HandlerFuncResource {
	return a.handlers
}
