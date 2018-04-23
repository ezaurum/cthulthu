package test

import (
	"fmt"
	"github.com/ezaurum/cthulthu/database"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/ezaurum/cthulthu/generators"
)

func DB(generators generators.IDGenerators) *gorm.DB {
	file := fmt.Sprintf("test%v.db", time.Now().UnixNano())
	db, _ := database.Open(generators, "sqlite3", file)
	return db
}
