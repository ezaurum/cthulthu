package context

type SessionConfig interface {

	SessionCookieName() string
	Domain() string
	SessionLifeLength() int
	PersistedCookieName() string
	SessionConfig(scn string, pcn string, domain string, len int)
}

func (a *app) SessionConfig(scn string, pcn string, domain string, len int) {
	a.sessionCookieName = scn
	a.persistedCookieName = pcn
	a.domain = domain
	a.sessionLifeLength = len
}

func (a *app) SessionCookieName() string {
	return a.sessionCookieName
}

func (a *app) Domain() string {
	return a.domain
}

func (a *app) SessionLifeLength() int {
	return a.sessionLifeLength
}

func (a *app) PersistedCookieName() string {
	return a.persistedCookieName
}
