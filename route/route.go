package route

import (
	"github.com/labstack/echo"
)

type Holder struct {
	RelativePath string
	Handler      echo.HandlerFunc
	Middleware   []echo.MiddlewareFunc
}

type Routes map[string][]Holder

func (routes Routes) POST(relativePath string, handlerFunc echo.HandlerFunc) Routes {
	return routes.Add("POST", relativePath, handlerFunc)
}

func (routes Routes) DELETE(relativePath string, handlerFunc echo.HandlerFunc) Routes {
	return routes.Add("DELETE", relativePath, handlerFunc)
}

func (routes Routes) PUT(relativePath string, handlerFunc echo.HandlerFunc) Routes {
	return routes.Add("PUT", relativePath, handlerFunc)
}

func (routes Routes) GET(relativePath string, handlerFunc echo.HandlerFunc) Routes {
	return routes.Add("GET", relativePath, handlerFunc)
}

func (routes Routes) HEAD(relativePath string, handlerFunc echo.HandlerFunc) Routes {
	return routes.Add("HEAD", relativePath, handlerFunc)
}

func (routes Routes) OPTIONS(relativePath string, handlerFunc echo.HandlerFunc) Routes {
	return routes.Add("OPTIONS", relativePath, handlerFunc)
}

func (routes Routes) Add(method string, relativePath string, handlerFunc echo.HandlerFunc) Routes {
	slice := routes[method]
	slice = append(slice, Holder{Handler: handlerFunc, RelativePath: relativePath})
	routes[method] = slice
	return routes
}

func AddAll(r *echo.Echo, routes Routes) {
	for k, v := range routes {
		switch k {
		case "GET":
			each(r.GET, v)
		case "POST":
			each(r.POST, v)
		case "DELETE":
			each(r.DELETE, v)
		case "PUT":
			each(r.PUT, v)
		case "OPTIONS":
			each(r.OPTIONS, v)
		case "HEAD":
			each(r.HEAD, v)
		}
	}
}

type routeFunc func(relativePath string, h echo.HandlerFunc,
	m ...echo.MiddlewareFunc) *echo.Route

func each(routeFunc routeFunc, holders []Holder) {
	for _, v := range holders {
		routeFunc(v.RelativePath, v.Handler, v.Middleware...)
	}
}

func InitRoute(r *echo.Echo, routes ...func() Routes) {
	for _, v := range routes {
		AddAll(r, v())
	}
}
