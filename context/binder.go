package context

type Binder interface {
	Bind(i interface{}) error
}

var _ Binder = &Request{}

func (r *Request) Bind(i interface{}) error {
	return r.Context().Bind(i)
}
