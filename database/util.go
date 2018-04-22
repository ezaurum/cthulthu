package database

import (
	"github.com/jinzhu/gorm"
)

func SaveAll(dbm *gorm.DB, targets ...interface{}) {
	transaction(dbm, func(tx *gorm.DB, v interface{}) {
		d := tx.Save(v)
		checkError(d, tx)
	}, targets...)
}
func Create(dbm *gorm.DB, target interface{}) interface{} {
	//TODO dbm.assignIDWhenNotAssigned(target)
	dbm.Create(target)
	return target
}

func CreateAll(db *gorm.DB, targets ...interface{}) {
	action := func(tx *gorm.DB, v interface{}) {
		d := tx.Create(v)
		checkError(d, tx)
	}
	transaction(db, action, targets...)
}

func checkError(d *gorm.DB, tx *gorm.DB) {
	if d.Error == nil {
		return
	}

	tx.Rollback()
	panic(d.Error)
}

func transaction(db *gorm.DB, action TransactionHandlerFunc, targets ...interface{}) {
	tx := db.Begin()
	for _, v := range targets {
		action(tx, v)
	}
	tx.Commit()
}

type TransactionHandlerFunc func(*gorm.DB, interface{})

func IsExist(db *gorm.DB, t interface{}, where ...interface{}) bool {
	dbR := db.Find(t, where...)
	switch dbR.Error {
	case nil:
		return true
	case gorm.ErrRecordNotFound:
		return false
	default:
		panic(db.Error)
	}
}
