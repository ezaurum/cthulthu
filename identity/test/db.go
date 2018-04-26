package test

import (
	"fmt"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/jinzhu/gorm"
	"time"
)

func DB(generators generators.IDGenerators) *gorm.DB {
	file := fmt.Sprintf("test%v.db", time.Now().UnixNano())
	db, _ := database.Open(generators, "sqlite3", file)
	return db
}
