package context

import (
	"net/http"
)

type Router interface {
	Handlers() []HandlerFuncResource
	GET(path string, handlerFunc RequestHandlerFunc, options ...interface{})
	POST(path string, handlerFunc RequestHandlerFunc, options ...interface{})
}

type router struct {
	handlers []HandlerFuncResource
}

var _ Router = &router{}

func (a *router) GET(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.handlers = append(a.handlers, HandlerFuncResource{
		Resource: Resource{
			Name:         "",
			ResourceType: HandlerFuncResourceType,
		},
		Method:      http.MethodGet,
		Path:        path,
		HandlerFunc: handlerFunc,
	})
}

func (a *router) POST(path string, handlerFunc RequestHandlerFunc, options ...interface{}) {
	a.handlers = append(a.handlers, HandlerFuncResource{
		Resource: Resource{
			Type:         DefaultHandler,
			Name:         "",
			ResourceType: HandlerFuncResourceType,
		},
		Method:      http.MethodPost,
		Path:        path,
		HandlerFunc: handlerFunc,
	})
}

func (a *router) Handlers() []HandlerFuncResource {
	return a.handlers
}
