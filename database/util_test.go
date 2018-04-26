package database

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/jinzhu/gorm"
	"testing"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/stretchr/testify/assert"
)

type Owner struct {
	Model
	Name string
	Pets []Pet
}

type Pet struct {
	Model
	Owner   Owner
	OwnerID int64
	Name    string
}

func TestIsExist(t *testing.T) {
	gens := snowflake.GetGenerators(0, &Owner{}, &Pet{})

	db := DB(gens)
	db.AutoMigrate(&Owner{}, &Pet{})

	owner := Owner{Name: "이름"}

	pets := []Pet{
		{
			Name: "멍멍1",
		},
		{
			Name: "멍멍2",
		},
	}
	owner.Pets = pets

	db.Create(&owner)

	assert.NotZero(t, owner.ID)
	assert.NotZero(t, pets[0].ID)
	assert.NotZero(t, pets[1].ID)
	assert.Equal(t, owner.ID, pets[0].OwnerID)
	assert.Equal(t, owner.ID, pets[1].OwnerID)

	var or Owner
	db.Preload("Pets").Find(&or, owner.ID)
	assert.Equal(t, len(owner.Pets), len(pets))

	var or1 Owner
	db.Find(&or1, owner.ID).Related(&or1.Pets, "OwnerID")
	assert.Equal(t, len(owner.Pets), len(or1.Pets))
}

func DB(generators generators.IDGenerators) *gorm.DB {
	//file := fmt.Sprintf("test%v.db", time.Now().UnixNano())
	db, _ := Open(generators, "sqlite3", "test.db")
	return db
}
