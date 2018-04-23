package database

import (
	"github.com/jinzhu/gorm"
	"github.com/ezaurum/cthulthu/generators"
	"reflect"
)

type DB struct {
	gorm.DB
	IDGenerators generators.IDGenerators
}

func (db *DB) AutoMigrate(values ...interface{}) {

	switch db.Dialect().GetName() {
	case "mysql":
		db.Set("gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;").
			AutoMigrate(values...)
		break
	default:
		db.DB.AutoMigrate(values...)
		break
	}
}

func (db *DB) Create(target interface{}) interface{} {
	db.assignIDWhenNotAssigned(target)
	db.DB.Create(target)
	return target
}

func (db *DB) Save(target interface{}) {
	db.DB.Save(target)
}

func (db *DB) SaveAll(targets ...interface{}) {
	db.transaction(func(tx *gorm.DB, v interface{}) {
		d := tx.Save(v)
		checkError(d, tx)
	}, targets...)
}

func (db *DB) CreateAll(targets ...interface{}) {
	action := func(tx *gorm.DB, v interface{}) {
		d := tx.Create(v)
		checkError(d, tx)
	}
	db.transaction(action, targets...)
}

type TransactionHandlerFunc func(*gorm.DB, interface{})

func (db *DB) transaction(action TransactionHandlerFunc, targets ...interface{}) {
	tx := db.Begin()
	for _, v := range targets {
		action(tx, v)
	}
}

func checkError(d *gorm.DB, tx *gorm.DB) {
	if d.Error == nil {
		return
	}

	tx.Rollback()
	panic(d.Error)
}

func (db *DB) assignIDWhenNotAssigned(target interface{}) int64 {
	stype := reflect.ValueOf(target).Elem()
	m := stype.FieldByName("Model")
	if m.Kind() == reflect.Struct {
		f0 := m.FieldByName("ID")
		if f0.IsValid() {
			if 0 == f0.Int() {
				id := db.IDGenerators[reflect.TypeOf(target).Name()].GenerateInt64()
				f0.SetInt(id)
			} else {
				return f0.Int()
			}
		}
	}
	return 0
}

func (db *DB) IsExist(t interface{}, where ...interface{}) bool {
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
