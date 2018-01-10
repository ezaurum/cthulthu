package identity

import "github.com/ezaurum/cthulthu/database"

type Identity struct {
	database.Model
	IdentityRole string
}

func (i Identity) Role() string {
	return i.IdentityRole
}
