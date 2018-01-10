package google

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ezaurum/cthulthu/session"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/identity"
	"net/url"
)


func TestGoogleRegister(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r := initializeTest(testDB, session.DefaultSessionExpires)

	form := getRegisterFormPostData()

	client := test.HttpClient{}
	w0 := client.PostFormRequest(r, "/google/register", form)

	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/", test.GetRedirect(w0))

	var result identity.OAuthIDToken
	b := testDB.IsExist(&result,
		&identity.OAuthIDToken{Provider: ProviderName,
		TokenID:form.Get("googleID"), Token:form.Get("googleToken") })
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

	form := getRegisterFormPostData()

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

	form := getRegisterFormPostData()

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
		AutoMigrates:            identity.AutoMigrates,
	})
}

func getRegisterFormPostData() url.Values {
	form := make(url.Values)
	form.Set("googleID", "112115032355210618122")
	form.Set("googleIDToken", "eyJhbGciOiJSUzI1NiIsImtpZCI6IjEwNWM2ZDVkNWIwYjI2ODA5ZmQxM2QxMzI5ZmJlY2E5ZmI0MTQyODIifQ.eyJhenAiOiI2Mjk4NzE3OTI3NjItdXZ0MTQxMDd1ajFzaGQzNWxxOWkwc2dvZHAyMHZkNzcuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2Mjk4NzE3OTI3NjItdXZ0MTQxMDd1ajFzaGQzNWxxOWkwc2dvZHAyMHZkNzcuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTIxMTUwMzIzNTUyMTA2MTgxMjIiLCJlbWFpbCI6ImV6YXVydW1AZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJ5LXdEUXhmc0ZNX2xudEk5LW5GN2pBIiwiaXNzIjoiYWNjb3VudHMuZ29vZ2xlLmNvbSIsImp0aSI6IjkzNTgzMDNlODc5YWJlZjY1MzkwYmE4NDIwMGY0N2E5NWFhZWU4ZTkiLCJpYXQiOjE1MTUzMTcyMzIsImV4cCI6MTUxNTMyMDgzMiwibmFtZSI6IuyhsOyEneq3nCIsInBpY3R1cmUiOiJodHRwczovL2xoNi5nb29nbGV1c2VyY29udGVudC5jb20vLUd4RW1UMFNrNVJZL0FBQUFBQUFBQUFJL0FBQUFBQUFBRUhFL1cxaUZpMV9UZ0dJL3M5Ni1jL3Bob3RvLmpwZyIsImdpdmVuX25hbWUiOiLshJ3qt5wiLCJmYW1pbHlfbmFtZSI6IuyhsCIsImxvY2FsZSI6ImVuIn0.jJEQZnKN2lJ2Mtz-xkD_XNImcb70gj1cpr7fircLldPZD12RKTP6jymfPN-rnd73mCG60kowMBJsu1Q0BSKhX_hVv3pYr5626TtjlOolS268i7CCNrAsVSw9pMnb9Py9MS2bvezte5bHt0lLECQcsHswJgm2Ehuy5PFORqgz3rPZrUcnSGk8LpTwLvCZAOq19LDhmNQ2kF6y6pRMWuwlu_MvBdyZ1p-dSKdpiVLn3NKxW8865JRsAVwH5LwACAkw5wMSDjQoLLwSmM81emTPyPCOdK3wN2FgqOYoeJvl5R9HEFi61AFOuR_ducSq05A7Q5VA9MVHuw0G054fY2YXzw")
	return form
}
