package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"time"
)

// 쿠키의 id 토큰 처리
func GetLoadCookieIDToken(dbm *database.Manager) authenticator.IDTokenLoader {
	return func(tokenString string) (authenticator.IDToken, bool) {
		var token CookieIDToken
		if dbm.IsExist(&token, &CookieIDToken{Token: tokenString}) {
			return token, true
		} else {
			return nil, false
		}
	}
}

func GetLoadIdentity(dbm *database.Manager) authenticator.IDLoader {
	return func(token authenticator.IDToken) (authenticator.Identity, bool) {
		identity := Identity{}

		if dbm.IsExist(&identity, token.IdentityKey()) {
			return identity, true
		}
		return nil, false
	}
}

func GetPersistToken(dbm *database.Manager) authenticator.TokenSaver {
	return func(token authenticator.IDToken) {
		dbm.Create(&CookieIDToken{
			IdentityID: token.IdentityKey(),
			Token:      token.TokenString(),
			Expires:    time.Now().Add(time.Hour * 24 * 365),
		})
	}
}