package identity

import (
	"testing"
	"github.com/ezaurum/cthulthu/database"
	"github.com/stretchr/testify/assert"
)

func TestOauthRegister(t *testing.T) {
	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()

	testDB.AutoMigrate(&Identity{}, &OAuthIDToken{})

	form := OAuthIDToken{
		TokenID:"test",
		Provider:"ProviderName",
		Token:"test token",
	}

	CreateIdentityByOAuth(form, testDB)

	var r OAuthIDToken
	b := testDB.IsExist(&r, &form)
	assert.True(t, b)
	assert.Equal(t,form.Token, r.Token)
	assert.Equal(t,form.TokenID, r.TokenID)
	assert.Equal(t,form.Provider, r.Provider)

	var i Identity
	testDB.IsExist(&i, r.IdentityKey())
}
