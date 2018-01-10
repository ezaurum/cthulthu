package identity

import (
	"github.com/ezaurum/cthulthu/database"
	"time"
)

func FindOAuthToken(r OAuthIDToken, dbm *database.Manager) (OAuthIDToken, bool) {
	var i OAuthIDToken
	exist := dbm.IsExist(&i, OAuthIDToken{
		TokenID:r.TokenID,
		Provider:r.Provider,
	})
	return i, exist
}

func UpdateOAuthToken(token OAuthIDToken, tokenString string, expires time.Time, dbm *database.Manager) {

	//TODO expires
	token.Token = tokenString
	dbm.Save(&token)


}
