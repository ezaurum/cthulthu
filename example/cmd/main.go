package main

import (
	"fmt"
	"github.com/ezaurum/cthulthu/app"
	"github.com/ezaurum/cthulthu/grant"
)

func main() {
	// Initialize the model from a string.
	text :=
		`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	line := `
p, alice, /alice_data/*, (GET)|(POST)
p, alice, /alice_data/resource1, POST
p, data_group_admin, /admin/*, POST
p, data_group_admin, /bob_data/*, POST
g, alice, data_group_admin
`
	if init, err := grant.Init(text, line); nil != err {
		panic(err)
	} else {
		fmt.Println(init)
	}

	app.Start("")
}
