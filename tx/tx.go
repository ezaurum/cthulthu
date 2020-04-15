package tx

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Request struct {
	Transaction *gorm.DB
	Error       error
}

func (rq *Request) Complete() error {
	if rq.Error != nil {
		return rq.Rollback()
	} else {
		return rq.Commit()
	}
}

func (rq *Request) Commit() error {
	if err := rq.Transaction.Commit(); nil != err.Error {
		return fmt.Errorf("commit error %w", err.Error)
	}
	return nil
}

func (rq *Request) Rollback() error {
	if err := rq.Transaction.Rollback(); nil != err.Error {
		return fmt.Errorf("rollback error %w", err.Error)
	}
	return nil
}

func New(db *gorm.DB) *Request {
	b := db.Begin()
	return &Request{
		Transaction: b,
		Error:       nil,
	}
}
