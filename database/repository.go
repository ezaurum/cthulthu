package database

import "github.com/jinzhu/gorm"

type Repository interface {
	Reader() *gorm.DB
	Writer() *gorm.DB
}

type DefaultRepository struct {
	masterDB *DBConnection
	slaveDB  *DBConnection
}

func (d *DefaultRepository) Reader() *gorm.DB {
	return d.slaveDB.DB
}

func (d *DefaultRepository) Writer() *gorm.DB {
	return d.masterDB.DB
}

func NewRepository(master *DBConnection, slave *DBConnection) Repository {
	return &DefaultRepository{
		masterDB: master,
		slaveDB:  slave,
	}
}
