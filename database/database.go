package database

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/jinzhu/gorm"
	"time"
)

type Model struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func Open(generators generators.IDGenerators,
	dialect string, args ...interface{}) (db *gorm.DB, err error) {
	db, err = gorm.Open(dialect, args...)
	if err != nil {
		return
	}

	switch dialect {
	case "mysql":
		db.Set("gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
	default:
	}

	db.Callback().Create().Before("gorm:create").
		Register("assign_int64_id", assignIDWhenNotAssigned(generators))

	return
}
