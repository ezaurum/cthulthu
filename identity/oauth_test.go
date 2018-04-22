package identity

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"reflect"
)

func TestOauthRegister(t *testing.T) {
	testDB := testDB()
	db := testDB.Connect()
	defer db.Close()

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

	CreateIdentityByOAuth(form, db, idGenerator)

	var r OAuthIDToken
	b := testDB.IsExist(&r, &form)
	assert.True(t, b)
	assert.Equal(t, form.Token, r.Token)
	assert.Equal(t, form.TokenID, r.TokenID)
	assert.Equal(t, form.Provider, r.Provider)

	var i Identity
	testDB.IsExist(&i, r.IdentityKey())
}
