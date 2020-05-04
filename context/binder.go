package context

import "github.com/ezaurum/cthulthu/errres"

type Binder interface {
	Bind(i interface{}) error
	Param(string) string
	QueryParam(name string) string
}

var _ Binder = &Request{}

func (r *Request) Bind(i interface{}) error {
	if err := r.Context().Bind(i); nil != err {
		return errres.BadReq("error bind", err, i)
	}
	return nil
}

func (r *Request) Param(name string) string {
	return r.Context().Param(name)
}

func (r *Request) QueryParam(name string) string {
	return r.Context().QueryParam(name)
}
