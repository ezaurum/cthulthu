package database

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"reflect"
)

func SaveAll(db *gorm.DB, targets ...interface{}) {
	Transaction(db, func(tx *gorm.DB, v interface{}) {
		d := tx.Save(v)
		checkError(d, tx)
	}, targets...)
}

type TransactionHandlerFunc func(*gorm.DB, interface{})

func Transaction(db *gorm.DB, action TransactionHandlerFunc, targets ...interface{}) {
	tx := db.Begin()
	for _, v := range targets {
		action(tx, v)
	}
	tx.Commit()
}

func checkError(d *gorm.DB, tx *gorm.DB) {
	if d.Error == nil {
		return
	}

	tx.Rollback()
	panic(d.Error)
}

func IsExist(db *gorm.DB, t interface{}, where ...interface{}) bool {
	dbR := db.Find(t, where...)
	switch dbR.Error {
	case nil:
		return true
	case gorm.ErrRecordNotFound:
		return false
	default:
		panic(dbR.Error)
	}
}

func assignIDWhenNotAssigned(generators generators.IDGenerators) func(scope *gorm.Scope) {
	return func(scope *gorm.Scope) {
		if scope.HasError() {
			return
		}

		primaryField := scope.PrimaryField()

		fieldType := primaryField.Field.Type().String()
		typeName := reflect.TypeOf(scope.Value).String()
		switch fieldType {
		case "int64":
			primaryField.Set(generators.GenerateInt64(typeName))
		case "string":
			primaryField.Set(generators.Generate(typeName))
		}
	}
}
