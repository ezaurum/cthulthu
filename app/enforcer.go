package app

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/ezaurum/cthulthu/app/adapter"
)

func initEnforcer() *casbin.Enforcer {
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
	m := model.NewModel()
	m.LoadModelFromText(text)

	line := `
p, alice, /alice_data/*, (GET)|(POST)
p, alice, /alice_data/resource1, POST
p, data_group_admin, /admin/*, POST
p, data_group_admin, /bob_data/*, POST
g, alice, data_group_admin
`
	sa := adapter.NewAdapter(line)
	e, _ := casbin.NewEnforcer(m)
	_ = e.InitWithModelAndAdapter(m, sa)
	return e
}
