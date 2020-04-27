package session

import (
	"errors"
	"fmt"
	"github.com/ezaurum/cthulthu/conv"
	"net/http"
	"time"
)

// 로그인 세션
type Session struct {
	ClientSession
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type ClientSession struct {
	ID          string `json:"id,omitempty"`
	AccountID   int64  `json:"account_id,string,omitempty"`
	Type        string `json:"type"`
	AccountName string `json:"account_name,omitempty"`
	Name        string `json:"name,omitempty"`
	MaxAge      int    `json:"max_age"`
}

func Populate(cs ClientSession) Session {
	now := time.Now()
	return Session{
		ClientSession: cs,
		CreatedAt:     now,
		UpdatedAt:     now,
		ExpiresAt:     now.Add(time.Second * time.Duration(cs.MaxAge)),
	}
}

func (s *Session) Extends() {
	now := time.Now()
	s.UpdatedAt = now
	s.ExpiresAt = now.Add(time.Second * time.Duration(s.MaxAge))
}

func PopulateAnonymous(maxAge int) (Session, error) {
	cs := ClientSession{
		Type:   AnonymousType,
		MaxAge: maxAge,
	}
	s := Populate(cs)

	return s, nil
}

func FromToken(token string) (*Session, error) {
	var s Session
	if err := conv.FromBase64Json(token, &s); nil != err {
		return &s, fmt.Errorf("error in convert to base64 %v, %w", s, err)
	}
	return &s, nil
}

func FromCookie(cookie *http.Cookie) (*Session, error) {
	var s Session
	if nil == cookie {
		return nil, ErrSessionEmpty
	}
	if err := conv.FromBase64Json(cookie.Value, &s); nil != err {
		return &s, fmt.Errorf("error in convert to base64 %v, %w", s, err)
	}
	now := time.Now()
	if s.ExpiresAt.Before(now) {
		return nil, ErrSessionExpired
	}
	return &s, nil
}
func (s *Session) ToCookie(httpCookieName string, cookieName string, domain string) (*http.Cookie, *http.Cookie, error) {
	if ck, err := s.ClientSession.ToCookie(cookieName, domain); nil != err {
		return nil, ck, err
	} else if sck, err := s.ToHttpOnlyCookie(httpCookieName); nil != err {
		return sck, ck, err
	} else {
		return sck, ck, nil
	}
}

func (s *Session) ToHttpOnlyCookie(cookieName string) (*http.Cookie, error) {
	json, err := conv.ToBase64Json(s)
	if nil != err {
		return nil, fmt.Errorf("error in marshaling session %v, %w", s, err)
	}
	return &http.Cookie{
		//todo 세션 타입에 따라 다르게 처리할 수 있도록,
		Name:     cookieName,
		Value:    json,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}, nil
}

func (cs *ClientSession) ToCookie(cookieName string, domain string) (*http.Cookie, error) {
	json, err := conv.ToBase64Json(cs)
	if nil != err {
		return nil, fmt.Errorf("error in marshaling client session %v, %w", cs, err)
	}
	return &http.Cookie{
		Name:     cookieName,
		Value:    json,
		Path:     "/",
		Domain:   domain,
		HttpOnly: false,
		MaxAge:   cs.MaxAge,
	}, nil
}

var ErrSessionRequired = errors.New("session required")
var ErrSessionEmpty = errors.New("session empty")
var ErrSessionExpired = errors.New("session expired")
