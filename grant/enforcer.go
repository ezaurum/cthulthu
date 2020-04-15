package grant

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/ezaurum/cthulthu/grant/adapter"
)

func Init(modelText, lineText string) (*casbin.Enforcer, error) {
	m := model.NewModel()
	if err := m.LoadModelFromText(modelText); nil != err {
		return nil, err
	}
	sa := adapter.NewAdapter(lineText)
	if e, err := casbin.NewEnforcer(); nil != err {
		return nil, err
	} else if err = e.InitWithModelAndAdapter(m, sa); nil != err {
		return e, err
	} else {
		return e, nil
	}
}
