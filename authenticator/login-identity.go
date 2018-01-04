package authenticator

import "time"

type LoginIdentity struct {
	id           int64
	UserID       string
	UserPassword string
	isPersisted  bool
	expires      time.Time
	Token        string
}

func (l LoginIdentity) TokenString() string {
	return l.Token
}

func (l LoginIdentity) IsPersisted() bool {
	return l.isPersisted
}

func (l LoginIdentity) IsExpired() bool {
	return time.Now().After(l.expires)
}
