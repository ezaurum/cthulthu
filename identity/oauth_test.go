package identity

import (
	"testing"
)

func TestOauthRegister(t *testing.T) {
	/*
		TODO
		targets := []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}}
		testDB := itest.DB()
		defer testDB.Close()

		testDB.AutoMigrate(&Identity{}, &OAuthIDToken{})

		form := OAuthIDToken{
			TokenID:  "test",
			Provider: "ProviderName",
			Token:    "test token",
		}

		gens := generators.New(func() generators.IDGenerator {
			return snowflake.New(0)
		}, &Identity{})
		name := reflect.TypeOf(&Identity{}).Name()
		idGenerator := gens[name]

		CreateIdentityByOAuth(form, testDB, idGenerator)

		var r OAuthIDToken
		b := database.IsExist(testDB,&r, &form)
		assert.True(t, b)
		assert.Equal(t, form.Token, r.Token)
		assert.Equal(t, form.TokenID, r.TokenID)
		assert.Equal(t, form.Provider, r.Provider)

		var i Identity
		database.IsExist(testDB, &i, r.IdentityKey()) */
}
