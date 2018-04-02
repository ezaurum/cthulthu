package route

import (
	"github.com/gin-gonic/gin"
)

type Holder struct {
	RelativePath string
	Handler      gin.HandlerFunc
}

type Routes map[string][]Holder

func (routes Routes) POST(relativePath string, handlerFunc gin.HandlerFunc) Routes {
	return routes.Add("POST", relativePath, handlerFunc)
}

func (routes Routes) GET(relativePath string, handlerFunc gin.HandlerFunc) Routes {
	return routes.Add("GET", relativePath, handlerFunc)
}

func (routes Routes) Add(method string, relativePath string, handlerFunc gin.HandlerFunc) Routes {
	slice := routes[method]
	slice = append(slice, Holder{Handler: handlerFunc, RelativePath: relativePath})
	routes[method] = slice
	return routes
}
func (routes Routes) AddPage(relativePath string, page string) Routes {
	return routes.Add("GET", relativePath, MakeJustHTML(page))
}

func (routes Routes) AddPageWith(relativePath string, page string, obj interface{}) Routes {
	return routes.Add("GET", relativePath, MakeHTMLWith(page, obj))
}

func AddTo(routes Routes, method string, relativePath string, handlerFunc gin.HandlerFunc) Routes {
	return routes.Add(method, relativePath, handlerFunc)
}

func AddAll(r *gin.Engine, routes Routes) {
	for k, v := range routes {
		switch k {
		case "GET":
			each(r.GET, v)
			break
		case "POST":
			each(r.POST, v)
			break
		case "DELETE":
			each(r.DELETE, v)
			break
		case "PUT":
			each(r.PUT, v)
			break
		}
	}
}

type routeFunc func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

func each(routeFunc routeFunc, holders []Holder) {
	for _, v := range holders {
		routeFunc(v.RelativePath, v.Handler)
	}
}

func InitRoute(r *gin.Engine, routes ...func() Routes) {
	for _, v := range routes {
		AddAll(r, v())
	}
}
