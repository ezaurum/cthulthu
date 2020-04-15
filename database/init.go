package database

import (
	"github.com/ezaurum/cthulthu/generators"
	"os"
	"strings"
)

func env(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func New(generators generators.IDGenerators, debug bool) (Repository, error) {
	masterCS := masterCS()
	masterCS.Debug = debug
	if _, err := masterCS.Open(); nil != err {
		return nil, err
	} else {
		masterCS.AssignGenerators(generators)
	}

	slaveCS := slaveCS()
	slaveCS.Debug = debug
	if _, err := slaveCS.Open(); nil != err {
		return nil, err
	}

	return NewRepository(masterCS, slaveCS), nil
}

func masterCS() *DBConnection {
	cs := NewConnection()
	if addr := env("SS_MDB_ADDR"); len(addr) > 0 {
		cs.Address = addr
	}
	if pass := env("SS_DB_PASS"); len(pass) > 0 {
		cs.Password = pass
	}
	if dbType := env("SS_DB_TYPE"); len(dbType) > 0 {
		cs.Type = dbType
	}
	if user := env("SS_DB_USER"); len(user) > 0 {
		cs.User = user
	}
	if db := env("SS_DB_DB"); len(db) > 0 {
		cs.Database = db
	}

	return cs
}

func slaveCS() *DBConnection {
	cs := NewConnection()
	if addr := env("SS_SDB_ADDR"); len(addr) > 0 {
		cs.Address = addr
	}
	if pass := env("SS_DB_PASS"); len(pass) > 0 {
		cs.Password = pass
	}
	if dbType := env("SS_DB_TYPE"); len(dbType) > 0 {
		cs.Type = dbType
	}
	if user := env("SS_DB_USER"); len(user) > 0 {
		cs.User = user
	}
	if db := env("SS_DB_DB"); len(db) > 0 {
		cs.Database = db
	}
	return cs
}
