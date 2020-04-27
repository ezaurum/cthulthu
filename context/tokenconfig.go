package context

type TokenConfig interface {
	TokenName() string
	TokenConfig(tn string)
}

func (a *app) TokenConfig(tn string) {
	a.tokenName = tn
}
func (a *app) TokenName() string {
	return a.tokenName
}
