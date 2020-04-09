package tx

import "github.com/jinzhu/gorm"

func Run(db *gorm.DB, run RunnerFunc) error {
	tx, rollback := PrepareRollback(db)
	defer rollback()
	if e := run(tx); nil != e {
		tx.Rollback()
		return e
	}

	if commit := tx.Commit(); nil != commit.Error {
		tx.Rollback()
		return commit.Error
	}
	return nil
}

func PrepareRollback(db *gorm.DB) (*gorm.DB, func()) {
	begin := db.Begin()
	return begin, func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			begin.Rollback()
			panic(p)
		}
	}
}

type RunnerFunc func(*gorm.DB) error
