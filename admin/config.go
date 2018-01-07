package admin

import "github.com/ezaurum/cthulthu/database"

const (
	staticDir   = "admin/static"
	templateDir = "admin/templates"
)

var (
	defaultConfig = []interface{}{"admin/model.conf", "admin/policy.csv"}
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

