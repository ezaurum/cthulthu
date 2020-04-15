package test

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func NewConnection(t *testing.T, dbName string) (*database.DBConnection, func()) {
	testDB := database.NewConnection()
	testDB.Type = "sqlite3"
	testDB.Database = dbName
	testDBF, errDB := testDB.Open()
	assert.NoError(t, errDB, "db open error")

	return testDB, func() {
		testDBF.Close()
		os.Remove(testDB.Database + ".db")
	}
}

func NewRepo(t *testing.T, dbName string) (database.Repository, func()) {
	testConn, closeFunc := NewConnection(t, "handler")
	repo := database.NewRepository(testConn, testConn)
	return repo, closeFunc
}
