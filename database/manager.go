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
	idGenerators     map[string]generators.IDGenerator
	nodeNumber       int64
	db               *gorm.DB
}

func Default() *Manager {
	// mysql Connect
	return New("root:example@tcp(127.0.0.1:3306)/dev?charset=utf8&parseTime=True&loc=Local",
		"mysql", 0)
}

func New(connectionString string, dialect string, nodeNumber int64) *Manager {
	return &Manager{
		connectionString: connectionString,
		dialect:          dialect,
		idGenerators:     make(map[string]generators.IDGenerator),
		nodeNumber:nodeNumber,
	}
}

func (dbm *Manager) DB() *gorm.DB {
	return dbm.db
}
func (dbm *Manager) Generate(typeName string) int64 {
	return dbm.idGenerators[typeName].GenerateInt64()
}

func (dbm *Manager) GenerateByType(v interface{}) int64 {
	return dbm.idGenerators[reflect.TypeOf(v).Name()].GenerateInt64()
}

func (dbm *Manager) AutoMigrate(values ...interface{}) {

	for _, v := range values {
		n := snowflake.New(dbm.nodeNumber)
		dbm.idGenerators[reflect.TypeOf(v).Name()] = n
	}

	switch dbm.dialect {
	case "mysql":
		dbm.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;").AutoMigrate(values...)
		break
	default:
		dbm.db.AutoMigrate(values...)
		break
	}
}

func (dbm *Manager) SetDialect(Dialect string) {
	dbm.dialect = Dialect
}

func (dbm *Manager) SetConnection(ConnectString string) {
	dbm.connectionString = ConnectString
}

func (dbm *Manager) Connect() *gorm.DB {
	if len(dbm.connectionString) < 1 {
		panic("Connect string is empty")
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

func TestNew() *Manager {
	file := fmt.Sprintf("test%v.db", time.Now().UnixNano())
	return New(file, "sqlite3", 0)
}

func Test() *Manager {
	return New("test.db", "sqlite3", 0)
}
