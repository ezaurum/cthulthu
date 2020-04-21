package context

import (
	"net/http"
)

func (a *router) AddHandler(pathString string, method string, handlerFunc ...RequestHandlerFunc) {
	joinedPath := a.JoinedPath(pathString)
	a.handlers = append(a.handlers, HandlerFuncResource{
		Resource: Resource{
			Name:         "",
			ResourceType: HandlerFuncResourceType,
		},
		Method:      method,
		Path:        joinedPath,
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
