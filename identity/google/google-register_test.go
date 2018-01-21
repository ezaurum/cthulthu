package google

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/identity"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestGoogleRegister(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r := initializeTest(testDB, session.DefaultSessionExpires)
	route.AddAll(r, Register())

	form := getGoogleRegisterForm()

	client := test.HttpClient{}
	w0 := client.PostFormRequest(r, "/google/register", form)

	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/", test.GetRedirect(w0))

	var result identity.OAuthIDToken
	b := testDB.IsExist(&result,
		&identity.OAuthIDToken{Provider: ProviderName,
			TokenID: form.Get("tokenID"), Token: form.Get("token")})

	var identity identity.Identity
	b0 := testDB.IsExist(&identity, result.IdentityKey())

	assert.True(t, b)
	assert.True(t, b0)
}

func TestGoogleAfterRegisterAuthenticated(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r := initializeTest(testDB, session.DefaultSessionExpires)
	route.AddAll(r, Register())

	form := getGoogleRegisterForm()

	client := test.HttpClient{}
	client.PostFormRequest(r, "/register", form)

	w := client.GetRequest(r, "/")
	assert.True(t, !test.IsRedirect(w))
}

func TestGoogleRegisterRememberLogin(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r := initializeTest(testDB, 1)
	route.AddAll(r, Register())

	form := getGoogleRegisterForm()

	client := test.HttpClient{}
	client.PostFormRequest(r, "/register", form)

	time.Sleep(time.Second)

	w := client.GetRequest(r, "/")
	assert.True(t, !test.IsRedirect(w))
}

func initializeTest(manager *database.Manager, expires int) *gin.Engine {
	return identity.Initialize(&identity.Config{
		DBManager:               manager,
		NodeNumber:              0,
		SessionExpiresInSeconds: expires,
		AutoMigrates:            identity.DefaultAutoMigrates,
	})
}

func getGoogleRegisterForm() url.Values {
	form := make(url.Values)
	form.Set("tokenID", "tokenID")
	form.Set("token", "tokenTemp")
	return form
}
