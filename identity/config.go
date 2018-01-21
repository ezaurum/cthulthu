package identity

import "github.com/ezaurum/cthulthu/database"

const (
	staticDir   = "identity/static"
	templateDir = "identity/templates"
)

var (
	DefaultConfig       = []interface{}{"identity/model.conf", "identity/policy.csv"}
	DefaultAutoMigrates = []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}, &OAuthIDToken{}, &Score{}}
)

type Config struct {
	DBManager               *database.Manager
	TemplateDir             string
	StaticDir               string
	NodeNumber              int64
	SessionExpiresInSeconds int
	AuthorizerConfig        []interface{}
	AutoMigrates            []interface{}
}
