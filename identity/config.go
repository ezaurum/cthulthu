package identity

import "github.com/ezaurum/cthulthu/database"

const (
	staticDir   = "admin/static"
	templateDir = "admin/templates"
)

var (
	defaultConfig = []interface{}{"admin/model.conf", "admin/policy.csv"}
	AutoMigrates = []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}, &OAuthIDToken{}}
)

type Config struct {
	DBManager *database.Manager
	TemplateDir string
	StaticDir string
	NodeNumber int64
	SessionExpiresInSeconds int
	AuthorizerConfig []interface{}
	AutoMigrates []interface{}
}

