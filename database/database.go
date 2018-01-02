package database

import "time"

type Model struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (m *Model) BeforeUpdate() (err error) {
	m.UpdatedAt = time.Now()
	return
}

func (m *Model) BeforeCreate() (err error) {
	m.CreatedAt = time.Now()
	return
}
