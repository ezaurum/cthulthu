package context

type Runner interface {
	Run(logic RequestHandlerFunc) error
	RunAll(logic []RequestHandlerFunc) error
	OnPanic()
}

var _ Runner = &Request{}

type RequestHandlerFunc func(c *Request) error

func (r *Request) Run(logic RequestHandlerFunc) error {
	return func() error {
		defer r.OnPanic()
		return logic(r)
	}()
}

func (r *Request) RunAll(logic []RequestHandlerFunc) error {
	return func() error {
		defer r.OnPanic()
		for _, handlerFunc := range logic {
			if err := handlerFunc(r); nil != err {
				r.Error = err
				return err
			}
		}
		return nil
	}()
}

func (r *Request) OnPanic() {
	var p interface{}
	if p = recover(); nil == p {
		return
	}
	// 진행중이던 트랜잭션이 있으면 롤백,
	// 에러는 이미 패닉이니 상관 안 함
	_ = r.RollbackTx()
	panic(p)
}
