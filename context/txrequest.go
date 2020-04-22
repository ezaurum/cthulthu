package context

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/tx"
	"github.com/jinzhu/gorm"
)

type txRequest struct {
	readDB      *gorm.DB
	writeDB     *gorm.DB
	Transaction *tx.Request
}

func (r *txRequest) SetTxError(err error) {
	if nil != r.Transaction {
		r.Transaction.Error = err
	}
}

func (r *txRequest) Reader() *gorm.DB {
	return r.readDB
}

func (r *txRequest) Writer() *gorm.DB {
	return r.writeDB
}

var _ TxRequest = &txRequest{}

type TxRequest interface {
	database.Repository
	StartTx(db *gorm.DB) *gorm.DB
	CompleteTx() error
	RollbackTx() error
	Tx() *gorm.DB
	TxError() error
	SetTxError(error)
}

func (r *txRequest) StartTx(db *gorm.DB) *gorm.DB {
	r.Transaction = tx.New(db)
	return r.Transaction.Transaction
}

func (r *txRequest) CompleteTx() error {
	if nil == r.Transaction {
		return nil
	}
	err := r.Transaction.Complete()
	r.Transaction = nil
	return err
}

func (r *txRequest) RollbackTx() error {
	if nil == r.Transaction {
		return nil
	}
	err := r.Transaction.Rollback()
	r.Transaction = nil
	return err
}

func (r *txRequest) RollbackOnPanic() {
	// 이전에 패닉이 있고 진행중이던 트랜잭션이 있으면
	if p := recover(); nil != p && r.Transaction != nil {
		// 이미 패닉이니 상관 안 함
		_ = r.Transaction.Rollback()
		panic(p)
	}
}

func (r *txRequest) TxError() error {
	if nil != r.Transaction && nil != r.Transaction.Error {
		return r.Transaction.Error
	}
	return nil
}

func (r *txRequest) Tx() *gorm.DB {
	if nil != r.Transaction {
		if r.Transaction.Error != nil {
			// 롤백 되었는데 모르고 Tx를 호출하면 nil을 반환한다
			// 구조상 에러처리를 따로 넣으면 번잡스러워진다
			return nil
		}
		return r.Transaction.Transaction
	} else {
		return r.StartTx(r.writeDB)
	}
}
