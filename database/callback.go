package database

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/jinzhu/gorm"
	"reflect"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func assignIDWhenNotAssigned(generators generators.IDGenerators) func(scope *gorm.Scope) {
	return func(scope *gorm.Scope) {
		if scope.HasError() {
			return
		}

		primaryField := scope.PrimaryField()

		if !primaryField.HasDefaultValue {
			return
		}

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

func RegisterAutoIDAssign(db2 *gorm.DB, generators generators.IDGenerators) {
	db2.Callback().Create().Before("gorm:create").
		Register("assign_int64_id", assignIDWhenNotAssigned(generators))
}
