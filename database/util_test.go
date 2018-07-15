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

func TestCreate(t *testing.T) {
	gens := snowflake.GetGenerators(0, &Owner{}, &Pet{})

	db := DB(gens)
	db.AutoMigrate(&Owner{}, &Pet{})

	defaultID := int64(1234)
	owner0 := Owner{
		Name: "test0",
		Model: Model{
			ID: defaultID,
		},
	}

	owner1 := Owner{
		Name: "test1",
	}

	//when
	db.Create(&owner0)
	db.Create(&owner1)

	//then
	assert.Equal(t, defaultID, owner0.ID)
	assert.NotZero(t, owner1.ID)
}

func TestFind(t *testing.T) {
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

func TestUpdate(t *testing.T) {
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

	owner.Pets[0].Name = "멍멍111"
	owner.Name = "WTF"

	db.Save(&owner)

	var or1 Owner
	db.Find(&or1, owner.ID).Related(&or1.Pets, "OwnerID")
	assert.Equal(t, "멍멍111", or1.Pets[0].Name)
	assert.Equal(t, owner.Pets[0].Name, or1.Pets[0].Name)
}

func TestDelete(t *testing.T) {
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

	db.Delete(&owner.Pets[1])
	owner.Pets = owner.Pets[:1]

	var or1 Owner
	db.Find(&or1, owner.ID).Related(&or1.Pets, "OwnerID")
	assert.Equal(t, len(owner.Pets), len(or1.Pets))
}

func DB(generators generators.IDGenerators) *gorm.DB {
	db2, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	RegisterAutoIDAssign(db2, generators)
	return db2
}
