package context

type HandlerFuncResource struct {
	Resource
	Method      string
	Path        string
	HandlerFunc []RequestHandlerFunc
}

type RequestHandlerFunc func(c *Request) error
