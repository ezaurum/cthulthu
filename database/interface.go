package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UniqueEntity interface {
	CreateIfNotExist(manager *gorm.DB) (interface{}, bool)
}

type Entity interface {
	Create(db *gorm.DB) (bool, error)
	Update(db *gorm.DB) (bool, error)
	Read(db *gorm.DB) (bool, error)
	ReadyBy(db *gorm.DB, where string, args ...interface{}) (bool, error)
	Delete(db *gorm.DB) (bool, error)
}
