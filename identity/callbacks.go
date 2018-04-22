package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"time"
	"github.com/jinzhu/gorm"
)

// 쿠키의 id 토큰 처리
func GetLoadCookieIDToken(dbm *gorm.DB) authenticator.IDTokenLoader {
	return func(tokenString string) (authenticator.IDToken, bool) {
		var token CookieIDToken
		if database.IsExist(dbm, &token, &CookieIDToken{Token: tokenString}) {
			return token, true
		} else {
			return nil, false
		}
	}
}

func GetLoadIdentityByCookie(dbm *gorm.DB) authenticator.IDLoader {
	return func(cookie authenticator.IDToken) (authenticator.Identity, bool) {
		identity := Identity{}
		if database.IsExist(dbm, &identity, cookie.IdentityKey()) {
			return identity, true
		}

		return nil, false
	}
}

func GetPersistToken(dbm *gorm.DB) authenticator.TokenSaver {
	return func(token authenticator.IDToken) authenticator.IDToken {
		return database.Create(dbm, &CookieIDToken{
			IdentityID: token.IdentityKey(),
			Token:      token.TokenString(),
			Expires:    time.Now().Add(time.Hour * 24 * 365),
		}).(authenticator.IDToken)
	}
}
