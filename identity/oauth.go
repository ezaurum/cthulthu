package identity

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/jinzhu/gorm"
	"time"
)

func FindOAuthToken(r OAuthIDToken, dbm *gorm.DB) (OAuthIDToken, bool) {
	var i OAuthIDToken
	exist := database.IsExist(dbm, &i, OAuthIDToken{
		TokenID:  r.TokenID,
		Provider: r.Provider,
	})
	return i, exist
}

func UpdateOAuthToken(token OAuthIDToken, tokenString string, expires time.Time, dbm *gorm.DB) {

	//TODO expires
	token.Token = tokenString
	dbm.Save(&token)

}
