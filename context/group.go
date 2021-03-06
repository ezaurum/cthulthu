package context

import (
	"path"
	"strings"
)

type RouteGroup interface {
	Group(path string) Router
	JoinedPath(pathString string) string
	AddGroupHandler(handlerFunc ...RequestHandlerFunc)
	GroupHandlers() []RequestHandlerFunc
}

func (a *router) AddGroupHandler(handlerFunc ...RequestHandlerFunc) {
	a.groupHandlers = append(a.groupHandlers, handlerFunc...)
}

func (a *router) JoinedPath(pathString string) string {
	var join string
	if strings.HasPrefix(pathString, ".") {
		join = a.basePath + pathString
	} else {
		join = path.Join(a.basePath, pathString)
	}
	if a.parent != nil {
		join = a.parent.JoinedPath(join)
	}
	return join
}

func (a *router) Group(path string) Router {
	r := &router{
		basePath: path,
		parent:   a,
	}
	a.children = append(a.children, r)
	return r
}

func (a *router) GroupHandlers() []RequestHandlerFunc {
	return a.groupHandlers
}
