package database

import (
	"github.com/jinzhu/gorm"
	"reflect"

	"github.com/ezaurum/cthulthu/generators"

	//TODO 임포트 자체를 바꿔야 하나?
	"fmt"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

// gorm 래퍼가 되지
type Manager struct {
	connectionString string
	dialect          string
	nodes            map[string]generators.IDGenerator
	NodeNumber       int64
	db               *gorm.DB
}

func Default() *Manager {
	// mysql connect
	return &Manager{
		connectionString: "root:example@tcp(127.0.0.1:3306)/dev?charset=utf8&parseTime=True&loc=Local",
		dialect:          "mysql",
		nodes:            make(map[string]generators.IDGenerator),
	}
}

func (dbm *Manager) DB() *gorm.DB {
	return dbm.db
}
func (dbm *Manager) Generate(typeName string) int64 {
	return dbm.nodes[typeName].GenerateInt64()
}

func (dbm *Manager) GenerateByType(v interface{}) int64 {
	return dbm.nodes[reflect.TypeOf(v).Name()].GenerateInt64()
}

func (dbm *Manager) AutoMigrate(values ...interface{}) {

	for _, v := range values {
		n := snowflake.New(dbm.NodeNumber)
		dbm.nodes[reflect.TypeOf(v).Name()] = n
	}

	switch dbm.dialect {
	case "mysql":
		dbm.db.Set("gorm:table_options", "ENGINE=InnoDB;CHARACTER SET=UTF-8")
		break
	}
	dbm.db.AutoMigrate(values...)
	node := snowflake.New(0)
	dbm.nodes["Spin"] = node
}

func (dbm *Manager) SetDialect(dialect string) {
	dbm.dialect = dialect
}

func (dbm *Manager) SetConnection(connectString string) {
	dbm.connectionString = connectString
}

func (dbm *Manager) Connect() *gorm.DB {
	if len(dbm.connectionString) < 1 {
		panic("connect string is empty")
	}

	db, err := gorm.Open(dbm.dialect, dbm.connectionString)
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	dbm.db = db
	return db
}
func (dbm *Manager) Create(target interface{}) interface{} {
	dbm.assignIDWhenNotAssigned(target)
	dbm.db.Create(target)
	return target
}

func (dbm *Manager) Save(target interface{}) {
	dbm.db.Save(target)
}

func (dbm *Manager) SaveAll(targets ...interface{}) {
	dbm.transaction(func(tx *gorm.DB, v interface{}) {
		d := tx.Save(v)
		checkError(d, tx)
	}, targets...)
}

func (dbm *Manager) CreateAll(targets ...interface{}) {
	action := func(tx *gorm.DB, v interface{}) {
		d := tx.Create(v)
		checkError(d, tx)
	}
	dbm.transaction(action, targets...)
}

func checkError(d *gorm.DB, tx *gorm.DB) {
	if d.Error == nil {
		return
	}

	tx.Rollback()
	panic(d.Error)
}

func (dbm *Manager) transaction(action TransactionHandlerFunc, targets ...interface{}) {
	tx := dbm.db.Begin()
	for _, v := range targets {
		action(tx, v)
	}
	tx.Commit()
}

type TransactionHandlerFunc func(*gorm.DB, interface{})

func (dbm *Manager) assignIDWhenNotAssigned(target interface{}) int64 {
	stype := reflect.ValueOf(target).Elem()
	m := stype.FieldByName("Model")
	if m.Kind() == reflect.Struct {
		f0 := m.FieldByName("ID")
		if f0.IsValid() {
			if 0 == f0.Int() {
				id := dbm.Generate(reflect.TypeOf(target).Name())
				f0.SetInt(id)
			} else {
				return f0.Int()
			}
		}
	}
	return 0
}

func (dbm *Manager) Find(token interface{}, where ...interface{}) interface{} {
	db := dbm.db.Find(token, where...)
	if nil != db.Error {
		return db.Error
	}
	return nil
}

func (dbm *Manager) IsExist(t interface{}, where ...interface{}) bool {
	error := dbm.Find(t, where...)
	switch error {
	case nil:
		return true
	case gorm.ErrRecordNotFound:
		return false
	default:
		panic(error)
	}
}

// 이거 테스트 때만 쓰긴 하는데...

func TestNew() *Manager {

	file := fmt.Sprintf("test%v.db", time.Now().UnixNano())

	// mysql connect
	return &Manager{
		connectionString: file,
		dialect:          "sqlite3",
		nodes:            make(map[string]generators.IDGenerator),
	}
}

func Test() *Manager {
	// mysql connect
	return &Manager{
		connectionString: "test.db",
		dialect:          "sqlite3",
		nodes:            make(map[string]generators.IDGenerator),
	}
}
