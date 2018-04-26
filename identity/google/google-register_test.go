package google

import (
	"net/url"
	"testing"
)

func TestGoogleRegister(t *testing.T) {
	/*
		db := itest.DB()
		defer db.Close()
		//TODO r := initializeTest(db, authenticator.DefaultSessionExpires)
		var r *gin.Engine
		route.AddAll(r, Register())

		form := getGoogleRegisterForm()

		client := test.HttpClient{}
		w0 := client.PostFormRequest(r, "/google/register", form)

		assert.True(t, test.IsRedirect(w0))
		assert.Equal(t, "/", test.GetRedirect(w0))

		var result identity.OAuthIDToken
		b := db.IsExist(&result,
			&identity.OAuthIDToken{Provider: ProviderName,
				TokenID: form.Get("tokenID"), Token: form.Get("token")})

		var identity identity.Identity
		b0 := db.IsExist(&identity, result.IdentityKey())

		assert.True(t, b)
		assert.True(t, b0)*/
}

func TestGoogleAfterRegisterAuthenticated(t *testing.T) {
	/*
		db := itest.DB()
		defer db.Close()
		//TODO r := initializeTest(db, authenticator.DefaultSessionExpires)
		var r *gin.Engine
		route.AddAll(r, Register())

		form := getGoogleRegisterForm()

		client := test.HttpClient{}
		client.PostFormRequest(r, "/register", form)

		w := client.GetRequest(r, "/")
		assert.True(t, !test.IsRedirect(w))*/
}

func TestGoogleRegisterRememberLogin(t *testing.T) {
	/*
		db := itest.DB()
		defer db.Close()
		//TODO r := initializeTest(db, 1)
		var r *gin.Engine
		route.AddAll(r, Register())

		form := getGoogleRegisterForm()

		client := test.HttpClient{}
		client.PostFormRequest(r, "/register", form)

		time.Sleep(time.Second)

		w := client.GetRequest(r, "/")
		assert.True(t, !test.IsRedirect(w))
	*/
}

func getGoogleRegisterForm() url.Values {
	form := make(url.Values)
	form.Set("tokenID", "tokenID")
	form.Set("token", "tokenTemp")
	return form
}
