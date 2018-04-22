package database

import (
	"github.com/jinzhu/gorm"
	"github.com/ezaurum/cthulthu/generators"
)

type DB struct {
	gorm.DB
	idGenerators generators.IDGenerators
}

func (db *DB) SaveAll(targets ...interface{}) {
	db.transaction(func(tx *gorm.DB, v interface{}) {
		d := tx.Save(v)
		checkError(d, tx)
	}, targets...)
}

func (db *DB) transaction(action TransactionHandlerFunc, targets ...interface{}) {
	tx := db.Begin()
	for _, v := range targets {
		action(tx, v)
	}
}
