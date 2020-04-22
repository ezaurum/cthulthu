package context

type Binder interface {
	Bind(i interface{}) error
	Param(string) string
	QueryParam(name string) string
}

var _ Binder = &Request{}

func (r *Request) Bind(i interface{}) error {
	return r.Context().Bind(i)
}

func (r *Request) Param(name string) string {
	return r.Context().Param(name)
}

func (r *Request) QueryParam(name string) string {
	return r.Context().QueryParam(name)
}
