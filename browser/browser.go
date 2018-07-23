package browser

import "github.com/ezaurum/cthulthu/database"

type Browser struct {
	database.Model
	Agent string
	IP string
}
