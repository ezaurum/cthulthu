package context

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Router interface {
	HandlerSetter
	RouteGroup
	Assign(e *echo.Echo, ctx Application, handlers ...RequestHandlerFunc) error
}

func (a *router) Assign(e *echo.Echo, ctx Application, parentHandlers ...RequestHandlerFunc) error {
	groupHandlers := append(parentHandlers, a.groupHandlers...)
	for _, handler := range a.handlers {
		handlerFuncs := append(groupHandlers, handler.HandlerFunc...)
		joinedPath := a.JoinedPath(handler.Path)
		defaultHandler := DefaultHandler(ctx, handlerFuncs...)

		switch handler.Method {
		case http.MethodGet:
			e.GET(joinedPath, defaultHandler)
		case http.MethodPost:
			e.POST(joinedPath, defaultHandler)
		case http.MethodPatch:
			e.PATCH(joinedPath, defaultHandler)
		case http.MethodPut:
			e.PUT(joinedPath, defaultHandler)
		case http.MethodDelete:
			e.DELETE(joinedPath, defaultHandler)
		case http.MethodConnect:
			e.CONNECT(joinedPath, defaultHandler)
		case http.MethodOptions:
			e.OPTIONS(joinedPath, defaultHandler)
		case http.MethodTrace:
			e.TRACE(joinedPath, defaultHandler)
		case http.MethodHead:
			e.HEAD(joinedPath, defaultHandler)
		}
	}

	for _, child := range a.children {
		_ = child.Assign(e, ctx, groupHandlers...)
	}

	return nil
}

type router struct {
	handlers      []HandlerFuncResource
	parent        Router
	basePath      string
	groupHandlers []RequestHandlerFunc
	children      []Router
}

var _ Router = &router{}
