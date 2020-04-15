package database

import (
	"fmt"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type DBConnection struct {
	User              string
	Password          string
	Database          string
	Address           string
	Port              int
	Location          string
	Type              string
	Debug             bool
	MaxOpenConnection int
	MaxIdleConnection int
	DB                *gorm.DB
	IDGenerators      generators.IDGenerators
}

func (cs *DBConnection) DSN() string {
	switch cs.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s", cs.User, cs.Password, cs.Address, cs.Port, cs.Database, cs.Location)
	case "sqlite3":
		return fmt.Sprintf("%s.db", cs.Database)
	default:
		return ""
	}
}

func (cs *DBConnection) ToString() string {
	switch cs.Type {
	case "mysql":
		return fmt.Sprintf("DSN: %s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s, Debug:%v, type:%v", cs.User, "*", cs.Address, cs.Port, cs.Database, cs.Location, cs.Debug, cs.Type)
	case "sqlite3":
		return fmt.Sprintf("DSN: %s.db, Debug:%v, type:%v", cs.Database, cs.Debug, cs.Type)
	default:
		return ""
	}
}

func (cs *DBConnection) Open() (*gorm.DB, error) {
	db, err := gorm.Open(cs.Type, cs.DSN())
	if nil != err && strings.Contains(err.Error(), "connect: connection refused") {
		time.Sleep(10 * time.Second)
		db, err = gorm.Open(cs.Type, cs.DSN())
	}
	if nil != err {
		return nil, err
	}
	switch cs.Type {
	case "mysql":
		db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
		db.DB().SetMaxIdleConns(cs.MaxIdleConnection)
		db.DB().SetMaxOpenConns(cs.MaxOpenConnection)
		db.DB().SetConnMaxLifetime(time.Hour)
	default:
	}
	if cs.Debug {
		db = db.Debug()
	}
	db.SingularTable(true)
	cs.DB = db
	return db, err
}

func (cs *DBConnection) AssignGenerators(idGenerators generators.IDGenerators) {
	cs.IDGenerators = idGenerators
	RegisterAutoIDAssign(cs.DB, cs.IDGenerators)
}

func NewConnection() *DBConnection {
	return &DBConnection{
		User:              "root",
		Password:          "devp",
		Database:          "dev",
		Address:           "localhost",
		Port:              3306,
		Location:          "Local",
		Type:              "mysql",
		Debug:             false,
		MaxOpenConnection: 50,
		MaxIdleConnection: 50,
	}
}
