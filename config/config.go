package config

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
	"html/template"
	"github.com/ezaurum/cthulthu/database"
)

type Config struct {
	DB           *gorm.DB
	AutoMigrates []interface{}

	NodeNumber              int64
	Generators              generators.IDGenerators
	SessionExpiresInSeconds int

	FuncMap template.FuncMap

	Address string

	Db  DBConfig
	Dir DirConfig
}

type DBConfig struct {
	Connection string
	Dialect    string
}

type DirConfig struct {
	Static   string
	Template string
}

func (cnf *Config) FromFile(configFile string) {
	toml, err := toml.LoadFile(configFile)
	if nil != err {
		panic(err)
	}
	toml.Unmarshal(cnf)
}
func (cnf *Config) InitDB() *gorm.DB {
	//Init DB
	db, err := gorm.Open(cnf.Db.Dialect, cnf.Db.Connection)
	if err != nil {
		panic(err)
	}
	database.RegisterAutoIDAssign(db, cnf.Generators)
	db.SingularTable(true)
	db.AutoMigrate(cnf.AutoMigrates...)

	cnf.DB = db
	return db
}
